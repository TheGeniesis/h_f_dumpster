# Cleanup branches script

This script allows to copy Gitlab CI CD variable under a new name

## Flags

- old - original variable name
- new - new variable name

## Setup

- Execute `task setup`
- Go to the [access token page](https://git.rxcorp.com/-/profile/personal_access_tokens)
- Generate access key with `api` priviledge
- Copy key and update `.env` file

## Start the application

- Execute `task execute OLD=<old-var> NEW=<new_var>`

## Development

### Prerequisites

- [Install staticcheck](https://staticcheck.io/docs/getting-started/#distribution-packages)

### Static analysis

- Run `task format` or `${GOPATH}/bin/staticcheck ./...`

[<< Go back](../README.md)
