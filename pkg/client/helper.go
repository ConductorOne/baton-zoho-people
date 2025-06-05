package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// By default, the number of objects returned per page is 100. The maximum number of object supported per page.
// It can be adjusted by adding the 'per_page' parameter in the query string.

const (
	ItemsPerPage   = 100
	accessTokenUrl = "https://accounts.zoho.%s/oauth/v2/token" // #nosec
)

type PageOptions struct {
	PageSize  int    `url:"limit,omitempty"`
	PageToken string `url:"sIndex,omitempty"`
}

var (
	TokenURLMap = map[string]string{
		"US": "com",
		"AU": "com.au",
		"EU": "eu",
		"IN": "in",
		"CN": "com.cn",
	}
)

func getPageSize(pageSize int) int {
	if pageSize <= 0 || pageSize > ItemsPerPage {
		pageSize = ItemsPerPage
	}
	return pageSize
}

func getNextPageToken(prevPageToken string, pageSize, recordsCount int) string {
	if prevPageToken == "" {
		return "1"
	}

	prevToken, _ := strconv.Atoi(prevPageToken)
	pageSize = getPageSize(pageSize)

	if recordsCount < pageSize {
		return ""
	}

	return strconv.Itoa(prevToken + pageSize)
}

type ReqOpt func(reqURL *url.URL)

func WithPageLimit(pageSize int) ReqOpt {
	return WithQueryParam("limit", strconv.Itoa(getPageSize(pageSize)))
}

func WithPageIndex(nextPageToken string) ReqOpt {
	if nextPageToken == "" {
		nextPageToken = "1"
	}
	return WithQueryParam("sIndex", nextPageToken)
}

func WithQueryParam(key string, value string) ReqOpt {
	return func(reqURL *url.URL) {
		q := reqURL.Query()
		q.Set(key, value)
		reqURL.RawQuery = q.Encode()
	}
}

func getTokenSource(
	ctx context.Context,
	clientId,
	clientSecret,
	refreshToken,
	domainAccount string,
) oauth2.TokenSource {
	cfg := clientcredentials.Config{
		EndpointParams: url.Values{
			"client_id":     []string{clientId},
			"client_secret": []string{clientSecret},
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refreshToken},
		},
		AuthStyle:    oauth2.AuthStyleInParams,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf(accessTokenUrl, TokenURLMap[domainAccount]),
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}
	return oauth2.ReuseTokenSource(token, cfg.TokenSource(ctx))
}
