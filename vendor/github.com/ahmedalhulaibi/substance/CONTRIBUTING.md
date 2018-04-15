# Contributing

This is my first open source project. For now the contributing guidelines will be loose and we'll develop it as we go.


## Contributing to the Contributing Guide

If there is a good practice that anyone wants to add to the contributing guide for the sustainability of the project then by all means open an issue and we will have an open discussion.


# Pull Requests

Make sure to have tests up to date. Please include any environmental setup instructions.

Currently I am the only reviewer. I will do my best to review within 1-2 weeks.

## Environment Setup

Install [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/)

Run the local hooks setup script if you want your code tested locally after each commit. The commits will be aborted if the `go test ./...` fails. The pre-commit script will also start the containers if required and export the required environment variable:
```bash
. setupScripts/setup-local-hooks.sh
```

Run [start-sql-containers.sh](https://github.com/ahmedalhulaibi/substance/blob/feature/gqlgen/setupScripts/setup-local-hooks.sh). This script will start postgres and mysql containers described in [docker-compose.yml](https://github.com/ahmedalhulaibi/substance/blob/feature/gqlgen/docker-compose.yml) and setup two environment variables
```bash
$ . setupScripts/start-sql-containers.sh
$ echo $SUBSTANCE_MYSQL 
root@tcp(172.19.0.3:3306)/delivery
$ echo $SUBSTANCE_PGSQL
postgres://travis_test:password@172.19.0.2:5432/postgres?sslmode=disable
```

# Issues

Open to any assistance in triaging and labelling issues. Hopefully we do not encounter too many problems/incidents.

## Bugs or Unexpected Behaviour

Please write the steps to reproduce an issue. Make sure to include which database you are using in the title of the 

## Enhancements

No restrictions as of this writing on feature requests. When requesting an enhancement some use cases would be nice.

