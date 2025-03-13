![Baton Logo](./baton-logo.png)

# `baton-zoho-people` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-zoho-people.svg)](https://pkg.go.dev/github.com/conductorone/baton-zoho-people) ![main ci](https://github.com/conductorone/baton-zoho-people/actions/workflows/main.yaml/badge.svg)

`baton-zoho-people` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Prerequisites
1. Use the [API Console](https://api-console.zoho.com) to create a Self Client
2. Use the Self Client to get the Client ID and Client Secret
3. Generate a code for one of the following scopes: `ZOHOPEOPLE.forms.ALL` or `ZOHOPEOPLE.forms.READ`

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-zoho-people
baton-zoho-people
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-zoho-people:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-zoho-people/cmd/baton-zoho-people@main

baton-zoho-people

baton resources
```

# Data Model

`baton-zoho-people` will pull down information about the following resources:
- Users
- Roles

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-zoho-people` Command Line Usage

```
baton-zoho-people

Usage:
  baton-zoho-people [flags]
  baton-zoho-people [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --domain-account               Domain-specific account region used to get access token (default "US")
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-zoho-people
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                 If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                      version for baton-zoho-people
      --zoho-client-id               (required) The Self Client zoho client id ($BATON_ZOHO_CLIENT_ID)
      --zoho-code                    (required) The authentication code generated using API Console ($BATON_ZOHO_CODE)
      --zoho-secret-id               (required) The Self Client zoho secret id ($BATON_ZOHO_SECRET_ID)

Use "baton-zoho-people [command] --help" for more information about a command.
```
