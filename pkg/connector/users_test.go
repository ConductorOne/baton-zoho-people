package connector

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-zoho-people/pkg/client"
	"github.com/conductorone/baton-zoho-people/test"
	"golang.org/x/oauth2"
)

var pageOptions = client.PageOptions{
	PageSize:  10,
	PageToken: "",
}

// Tests that the client can fetch users based on the documented API below.
// https://www.zoho.com/people/api/forms-api/fetch-single-section.html
func TestZohoPeopleClient_GetUsers(t *testing.T) {
	// Create a mock response.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(test.ReadFile("employeesMock.json"))),
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
