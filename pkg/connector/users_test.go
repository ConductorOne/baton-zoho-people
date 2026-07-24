package connector

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-zoho-people/pkg/client"
	"github.com/conductorone/baton-zoho-people/test"
	"golang.org/x/oauth2"

	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

var pageOptions = client.PageOptions{
	PageSize:  10,
	PageToken: "",
}

// Tests that the client can fetch users based on the documented API below.
// https://www.zoho.com/people/api/forms-api/fetch-single-section.html
func TestZohoPeopleClient_GetUsers(t *testing.T) {
	// Create a mock response.
	mockData, err := test.ReadFile("employeesMock.json")
	if err != nil {
		t.Fatal(err)
	}
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(mockData)),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// Create a test client with the mock response.
	testClient := test.NewTestClient(mockResponse, nil)

	// Call GetUsers
	ctx := context.Background()
	result, _, nextOptions, err := testClient.ListUsers(ctx, pageOptions)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the result.
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check count.
	if len(result) != 2 {
		t.Errorf("Expected Count to be 2, got %d", len(result))
	}

	for index, user := range result {
		expectedUser := client.Employee{
			ZohoID:     int64(index + 10000000000),
			FirstName:  test.Employees[index]["firstName"],
			LastName:   test.Employees[index]["lastName"],
			EmployeeID: test.Employees[index]["employeeID"],
			EmailID:    test.Employees[index]["email"],
		}

		if !reflect.DeepEqual(user.ZohoID, expectedUser.ZohoID) &&
			!reflect.DeepEqual(user.EmployeeID, expectedUser.EmployeeID) &&
			!reflect.DeepEqual(user.EmailID, expectedUser.EmailID) &&
			!reflect.DeepEqual(user.LastName, expectedUser.LastName) &&
			!reflect.DeepEqual(user.FirstName, expectedUser.FirstName) {
			t.Errorf("Unexpected user: got %+v, want %+v", user, expectedUser)
		}
	}

	// Check next options.
	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

func TestZohoPeopleClient_GetUsers_RequestDetails(t *testing.T) {
	// Create a custom RoundTripper to capture the request.
	var capturedRequest *http.Request
	mockTransport := &test.MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
			Header:     make(http.Header),
		},
		Err: nil,
	}
	mockTransport.Response.Header.Set("Content-Type", "application/json")

	mockRoundTrip := func(req *http.Request) (*http.Response, error) {
		capturedRequest = req
		return mockTransport.Response, mockTransport.Err
	}
	mockTransport.SetRoundTrip(mockRoundTrip)

	// Create a test client with the mock transport.
	httpClient := &http.Client{Transport: mockTransport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)

	token := oauth2.Token{
		AccessToken: "access-token-hash",
	}
	testClient := client.NewClient(oauth2.StaticTokenSource(&token), baseHttpClient)

	// Call GetUsers.
	ctx := context.Background()
	_, _, err, _ := testClient.ListUsers(ctx, pageOptions)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the request details.
	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}

	// Check URL components.
	expectedURL := "https://people.zoho.com/people/api/forms/employee/getRecords?limit=10&sIndex=1"
	if capturedRequest.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, capturedRequest.URL.String())
	}

	// Check headers.
	expectedHeaders := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": "Bearer access-token-hash",
	}

	for key, expectedValue := range expectedHeaders {
		if value := capturedRequest.Header.Get(key); value != expectedValue {
			t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, value)
		}
	}
}

// singleEmployeeMock is a minimal GetEmployeeByID response body with a
// non-empty Role, used to exercise the cross-type role-grant emission in
// userBuilder.Grants.
const singleEmployeeMock = `{
	"response": {
		"result": [
			{
				"EmployeeID": "S20",
				"EmailID": "christopherbrown@zylker.com",
				"FirstName": "Christopher",
				"LastName": "Brown",
				"Role": "Team Incharge",
				"Zoho_ID": 100000000000
			}
		],
		"message": "success",
		"status": 200
	}
}`

func newTestUserResource(t *testing.T, id string) *v2.Resource {
	t.Helper()
	return &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     id,
		},
	}
}

// TestUserBuilder_Grants_RoleNotSkipped verifies that when the "role"
// resource type is not excluded from the sync filter (skipRoleGrants=false,
// the zero-value default, or no filter is set), Grants emits the cross-type
// role-assignment grant sourced from the employee's Role field.
func TestUserBuilder_Grants_RoleNotSkipped(t *testing.T) {
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(singleEmployeeMock)),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	c := test.NewTestClient(mockResponse, nil)
	u := newUserBuilder(c, false)

	grants, _, err := u.Grants(context.Background(), newTestUserResource(t, "100000000000"), rs.SyncOpAttrs{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(grants) != 1 {
		t.Fatalf("expected 1 grant when skipRoleGrants=false, got %d", len(grants))
	}

	gotRoleID := grants[0].Entitlement.Resource.Id.Resource
	wantRoleID := GetRoleID("Team Incharge")
	if gotRoleID != wantRoleID {
		t.Errorf("expected grant against role %q, got %q", wantRoleID, gotRoleID)
	}
	if grants[0].Entitlement.Resource.Id.ResourceType != roleResourceType.Id {
		t.Errorf("expected grant resource type %q, got %q", roleResourceType.Id, grants[0].Entitlement.Resource.Id.ResourceType)
	}
}

// TestUserBuilder_ResourceType_SkipAnnotations verifies that ResourceType
// annotates the user resource type based on skipRoleGrants instead of Grants
// guarding internally: SkipEntitlements when the "role" resource type is not
// excluded from the sync filter (skipRoleGrants=false, the zero-value
// default - Grants must still run to emit the cross-type role-assignment
// grant), SkipEntitlementsAndGrants when it is excluded (skipRoleGrants=true
// - the SDK must skip Grants entirely, since the connector should never emit
// a grant referencing a resource type it isn't syncing). It also verifies
// the shared userResourceType package var is never mutated by either call.
func TestUserBuilder_ResourceType_SkipAnnotations(t *testing.T) {
	c := &client.ZohoPeopleClient{}

	roleNotSkipped := newUserBuilder(c, false)
	rtNotSkipped := roleNotSkipped.ResourceType(context.Background())
	notSkippedAnnos := annotations.Annotations(rtNotSkipped.Annotations)
	if !notSkippedAnnos.Contains(&v2.SkipEntitlements{}) {
		t.Error("expected skipRoleGrants=false ResourceType to contain SkipEntitlements")
	}
	if notSkippedAnnos.Contains(&v2.SkipEntitlementsAndGrants{}) {
		t.Error("expected skipRoleGrants=false ResourceType to NOT contain SkipEntitlementsAndGrants")
	}

	roleSkipped := newUserBuilder(c, true)
	rtSkipped := roleSkipped.ResourceType(context.Background())
	skippedAnnos := annotations.Annotations(rtSkipped.Annotations)
	if !skippedAnnos.Contains(&v2.SkipEntitlementsAndGrants{}) {
		t.Error("expected skipRoleGrants=true ResourceType to contain SkipEntitlementsAndGrants")
	}

	if len(userResourceType.Annotations) != 0 {
		t.Fatalf("expected shared userResourceType.Annotations to remain unmutated, got %d annotations", len(userResourceType.Annotations))
	}
}
