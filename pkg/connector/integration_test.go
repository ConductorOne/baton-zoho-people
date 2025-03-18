package connector

import (
	"context"
	"fmt"
	"os"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-zoho-people/pkg/client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	ctx              = context.Background()
	domainAccount, _ = os.LookupEnv("ZOHO_PEOPLE_DOMAIN_ACCOUNT")
	apiToken, _      = os.LookupEnv("ZOHO_PEOPLE_API_TOKEN")
	parentResourceID = &v2.ResourceId{}
	pToken           = &pagination.Token{Size: 50, Token: ""}
)

func initClient(t *testing.T) *client.ZohoPeopleClient {
	if apiToken == "" {
		message :=
			fmt.Sprintf("Any of the required params not found. Api token: %s", apiToken)
		t.Skip(message)
	}

	if domainAccount == "" {
		domainAccount = "US"
	}

	token := oauth2.Token{
		AccessToken: apiToken,
	}
	c, err := client.New(
		ctx,
		client.ZohoAuthData{},
		oauth2.StaticTokenSource(&token),
	)
	if err != nil {
		t.Errorf("ERROR: Failed to create client: %v", err)
	}
	return c
}

func TestUserBuilderList(t *testing.T) {
	c := initClient(t)

	u := newUserBuilder(c)
	res, _, _, err := u.List(ctx, parentResourceID, pToken)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	message := fmt.Sprintf("Amount of users obtained: %d", len(res))
	t.Log(message)
}

func TestRoleBuilderList(t *testing.T) {
	c := initClient(t)

	r := newRoleBuilder(c)

	res, _, _, err := r.List(ctx, parentResourceID, pToken)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	message := fmt.Sprintf("Amount of roles obtained: %d", len(res))
	t.Log(message)
}
