package test

import (
	"log"
	"net/http"
	"os"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-zoho-people/pkg/client"
	"golang.org/x/oauth2"
)

var (
	Employees = []map[string]string{
		{
			"firstName":  "Christopher",
			"lastName":   "Brown",
			"email":      "christopherbrown@zylker.com",
			"employeeID": "S20",
		},
		{
			"firstName":  "David",
			"lastName":   "Johnson",
			"email":      "michaeljohnson@zylker.com",
			"employeeID": "S19",
		}}
)

// Custom RoundTripper for testing.
type TestRoundTripper struct {
	response *http.Response
	err      error
}

type MockRoundTripper struct {
	Response  *http.Response
	Err       error
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func (m *MockRoundTripper) SetRoundTrip(roundTrip func(*http.Request) (*http.Response, error)) {
	m.roundTrip = roundTrip
}

func (t *TestRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return t.response, t.err
}

// Helper function to create a test client with custom transport.
func NewTestClient(response *http.Response, err error) *client.ZohoPeopleClient {
	transport := &TestRoundTripper{response: response, err: err}
	httpClient := &http.Client{Transport: transport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)

	token := oauth2.Token{
		AccessToken: "",
	}
	return client.NewClient(oauth2.StaticTokenSource(&token), baseHttpClient)
}

func ReadFile(fileName string) string {
	data, err := os.ReadFile("../../test/mockResponses/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
