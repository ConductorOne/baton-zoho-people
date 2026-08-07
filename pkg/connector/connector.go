package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-zoho-people/pkg/client"
	cfg "github.com/conductorone/baton-zoho-people/pkg/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Connector struct {
	client *client.ZohoPeopleClient
	// skipRoleResourceType reports whether role is excluded from the sync
	// filter. Named for the skip condition so the zero value is safe: main.go
	// registers a zero-value Connector{} as the capabilities factory.
	skipRoleResourceType bool
}

type Option func(*Connector) error

func (d *Connector) SetTokenSource(tokenSource oauth2.TokenSource) {
	d.client.TokenSource = tokenSource
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client, d.skipRoleResourceType),
		newRoleBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Zoho people Connector",
		Description: "Connector to sync and manage users and roles in Zoho people.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, zohoClientID, zohoSecretID, zohoRefreshToken, domainAccount string, skipRoleResourceType bool) (*Connector, error) {
	l := ctxzap.Extract(ctx)

	zohoPeopleClient, err := client.New(ctx, client.ZohoAuthData{
		ClientID:      zohoClientID,
		ClientSecret:  zohoSecretID,
		ClientCode:    zohoRefreshToken,
		DomainAccount: domainAccount,
	})
	if err != nil {
		l.Error("error creating Zoho People client", zap.Error(err))
		return nil, err
	}

	return &Connector{
		client:               zohoPeopleClient,
		skipRoleResourceType: skipRoleResourceType,
	}, nil
}

// NewLambdaConnector satisfies cli.NewConnector for use with config.RunConnector.
func NewLambdaConnector(ctx context.Context, ac *cfg.ZohoPeople, opts *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	// nil opts means no filter, so nothing is skipped.
	skipRoleResourceType := opts != nil && !opts.WillSyncResourceType(RoleResourceTypeID)

	cb, err := New(ctx, ac.ZohoClientId, ac.ZohoSecretId, ac.ZohoRefreshToken, ac.DomainAccount, skipRoleResourceType)
	if err != nil {
		return nil, nil, err
	}
	return cb, nil, nil
}
