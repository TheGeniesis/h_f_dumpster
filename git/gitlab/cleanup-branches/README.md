# Cleanup branches script

This script allows to remove stale git branches.

## Rules

- Branch is not protected
- Branch doesn't have MR
- Branch wasn't updated in last 2 months

## Flags

- dry-run - Execute code without removing branches
- scope - (projects/group) set scope to execute (get repositories to cleanup)

## Setup

- Execute `task setup`
- Go to the [access token page](https://git.rxcorp.com/-/profile/personal_access_tokens)
- Generate access key with `api` priviledge
- Copy key and update `.env` file

## Start the application

- Execute `task execute`

## Development

### Prerequisites

- [Install staticcheck](https://staticcheck.io/docs/getting-started/#distribution-packages)

### Static analysis

- Run `task format` or `${GOPATH}/bin/staticcheck ./...`

[<< Go back](../README.md)
