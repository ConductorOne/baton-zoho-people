package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/ratelimit"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/oauth2"
)

type ZohoPeopleClient struct {
	wrapper     *uhttp.BaseHttpClient
	TokenSource oauth2.TokenSource
}

type ZohoAuthData struct {
	ClientID      string
	ClientSecret  string
	ClientCode    string
	DomainAccount string
}

type Option func(client *ZohoPeopleClient)

const (
	baseUrl        = "https://people.zoho.com/people/api/forms"
	accessTokenUrl = "https://accounts.zoho.%s/oauth/v2/token" // #nosec

	getDepartmentRecords    = "/department/getRecords"
	getDepartmentByRecordId = "/department/getDataByID"
	getEmployeeRecords      = "/employee/getRecords"
	getEmployeeByRecordId   = "/employee/getDataByID"
)

func New(ctx context.Context, authData ZohoAuthData, authToken ...oauth2.TokenSource) (*ZohoPeopleClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, err
	}

	client := ZohoPeopleClient{
		wrapper: cli,
	}

	if authToken != nil {
		client.TokenSource = authToken[0]
	} else {
		client.TokenSource = getTokenSource(ctx, authData.ClientID, authData.ClientSecret, authData.ClientCode, authData.DomainAccount)
	}

	return &client, nil
}

func NewClient(tokenSource oauth2.TokenSource, httpClient ...*uhttp.BaseHttpClient) *ZohoPeopleClient {
	var wrapper = &uhttp.BaseHttpClient{}
	if httpClient != nil || len(httpClient) != 0 {
		wrapper = httpClient[0]
	}
	return &ZohoPeopleClient{
		wrapper:     wrapper,
		TokenSource: tokenSource,
	}
}

func (c *ZohoPeopleClient) ListUsers(ctx context.Context, options PageOptions) ([]Employee, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res EmployeeResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getEmployeeRecords)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating url: %s", err))
		return nil, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithPageIndex(options.PageToken), WithPageLimit(options.PageSize))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, "", nil, err
	}

	result := res.Response.Result
	if result != nil {
		var employees []Employee
		for _, item := range result {
			for _, empList := range item {
				employees = append(employees, empList...)
			}
		}
		return employees, getNextPageToken(options.PageToken, options.PageSize, len(result)), annotation, nil
	} else {
		return nil, "", annotation, nil
	}
}

func (c *ZohoPeopleClient) ListDepartments(ctx context.Context, options PageOptions) ([]Department, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res DepartmentResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getDepartmentRecords)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating url: %s", err))
		return nil, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithPageIndex(options.PageToken), WithPageLimit(options.PageSize))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, "", nil, err
	}

	result := res.Response.Result
	if result != nil {
		var departments []Department
		for _, item := range result {
			for _, empList := range item {
				departments = append(departments, empList...)
			}
		}
		return departments, getNextPageToken(options.PageToken, options.PageSize, len(result)), annotation, nil
	} else {
		return nil, "", annotation, nil
	}
}

func (c *ZohoPeopleClient) GetDepartmentByID(ctx context.Context, departmentID string) ([]Department, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res SingleDepartmentResponse
	var departments []Department
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getDepartmentByRecordId)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating url: %s", err))
		return departments, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithQueryParam("recordId", departmentID))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resource: %s", err))
		return departments, "", nil, err
	}

	result := res.Response.Result
	if result != nil {
		departments = append(departments, result...)
	}
	return departments, "", annotation, nil
}

func (c *ZohoPeopleClient) GetEmployeeByID(ctx context.Context, employeeID string) ([]Employee, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res SingleEmployeeResponse
	var employees []Employee
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getEmployeeByRecordId)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating url: %s", err))
		return employees, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithQueryParam("recordId", employeeID))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resource: %s", err))
		return employees, "", nil, err
	}

	result := res.Response.Result
	if result != nil {
		employees = append(employees, result...)
	}
	return employees, "", annotation, nil
}

func (c *ZohoPeopleClient) getResourcesFromAPI(
	ctx context.Context,
	urlAddress string,
	res any,
	reqOptions ...ReqOpt,
) (annotations.Annotations, error) {
	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, reqOptions...)

	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (c *ZohoPeopleClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	reqOptions ...ReqOpt,
) (http.Header, annotations.Annotations, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(endpointUrl)

	if err != nil {
		return nil, nil, err
	}

	for _, o := range reqOptions {
		o(urlAddress)
	}

	authToken, err := c.TokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	req, err := c.wrapper.NewRequest(
		ctx,
		method,
		urlAddress,
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithAcceptJSONHeader(),
	)
	authToken.SetAuthHeader(req)

	if err != nil {
		return nil, nil, err
	}

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost:
		var doOptions []uhttp.DoOption
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		resp, err = c.wrapper.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}
	case http.MethodDelete:
		resp, err = c.wrapper.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	annotation := annotations.Annotations{}
	if resp != nil {
		if desc, err := ratelimit.ExtractRateLimitData(resp.StatusCode, &resp.Header); err == nil {
			annotation.WithRateLimiting(desc)
		} else {
			return nil, annotation, err
		}

		return resp.Header, annotation, nil
	}

	return nil, nil, err
}
