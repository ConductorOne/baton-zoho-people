package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	clientIDField = field.StringField(
		"zoho-client-id",
		field.WithRequired(true),
		field.WithDisplayName("Zoho client ID"),
		field.WithDescription("Client ID of the Self Client Application for Zoho."),
		field.WithPlaceholder("Your Zoho client ID"),
	)
	secretIDField = field.StringField(
		"zoho-secret-id",
		field.WithRequired(true),
		field.WithDisplayName("Zoho secret ID"),
		field.WithDescription("Secret ID of the Self Client Application for Zoho."),
		field.WithPlaceholder("Your Zoho secret ID"),
		field.WithIsSecret(true),
	)
	refreshTokenField = field.StringField(
		"zoho-refresh-token",
		field.WithRequired(true),
		field.WithDisplayName("Zoho refresh token"),
		field.WithDescription("The temporary authorization code to access Zoho APIs."),
		field.WithPlaceholder("Your Zoho refresh token"),
		field.WithIsSecret(true),
	)
	domainAccount = field.SelectField(
		"domain-account",
		[]string{"US", "AU", "EU", "IN", "CN"},
		field.WithDescription("The domain specific account to get the access token."),
		field.WithDisplayName("Domain account"),
		field.WithDefaultValue("US"),
	)
)

var Config = field.NewConfiguration(
	[]field.SchemaField{clientIDField, secretIDField, refreshTokenField, domainAccount},
	field.WithConnectorDisplayName("Zoho People"),
	field.WithIconUrl("/static/app-icons/zoho-people.svg"),
	field.WithHelpUrl("/docs/baton/zoho-people"),
)
