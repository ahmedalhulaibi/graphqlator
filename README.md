# Graphqlator CLI
__*This project is a WIP.*__

Graphqlator takes your existing database schema and generates code for a GraphQL-Go server. Type 'graphqlator help' to see usage.

## Status

| Build                             | Report card                           |
| :-------------------------------: | :-----------------------------------: |
| [![Build][build-badge]][build-link] | [![Report card][rc-badge]][rc-link]   |

[build-badge]: https://travis-ci.org/ahmedalhulaibi/graphqlator.svg?branch=master "Travis-CI build status"
[build-link]: https://travis-ci.org/ahmedalhulaibi/graphqlator "Travis-CI build status link"
[rc-badge]: https://goreportcard.com/badge/github.com/ahmedalhulaibi/graphqlator "Report card status"
[rc-link]: https://goreportcard.com/report/github.com/ahmedalhulaibi/graphqlator "Report card status"

## Supported Data Stores:

- mysql
- mariadb
- postgres

## Installation:

```
go get github.com/ahmedalhulaibi/graphqlator
```

Or Download prebuilt binaries from the [releases page](https://github.com/ahmedalhulaibi/go-graphqlator-cli/releases)

## Prerequisites

[goreturns](https://github.com/sqs/goreturns) - Generator uses goreturns to remove unnecessary generated imports

[grahpql-go](https://github.com/graphql-go/graphql) - Generated code uses graphql-go

[GORM](https://github.com/jinzhu/gorm) - Generated code uses GORM

## Usage
```
  graphqlator [flags]
  graphqlator [command]
```
Available Commands:
```
  init        Create a graphqlator-pkg.json file.
  describe    Describe database or table
  generate    Generate GraphQL type schema from database table.
  help        Help about any command
  version     Print the version number of Graphqlator
```
Flags:
  -h, --help   help for graphqlator

Use "graphqlator [command] --help" for more information about a command.

## Example Usage:

### graphqlator init

This command will walk you through the creation of a graphqlator-pkg.json file. This file can be edited after the fact.

```
graphqlator init
Input project name (enter to continue): first-graphql-server
Input database type e.g. mysql (enter to continue): mysql
Input Database Connection String
				MySql Example: username:password@tcp(localhost:3306)/schemaname
				Postgresql Example: postgres://username:password@localhost:5432/dbname
Input database connection string (enter to continue): ahmed:password@tcp(localhost:3306)/delivery
Input git repo url (enter to continue): 
Input table names - Must be EXACT spelling and case (enter without input to skip)
Table #1 : Persons
Table #2 : Orders
Table #3 : AntiOrders
Table #4 : 
```
#### graphqlator-pkg.json:
```
{
    "project_name": "first-graphql-server",
    "database_type": "mysql",
    "connection_string": "ahmed:password@tcp(localhost:3306)/delivery",
    "table_names": [
        "Persons",
        "Orders",
        "AntiOrders"
    ],
    "git_repo": ""
}
```
### graphqlator generate
This command will use the information in graphqlator-pkg.json and generate two go files: 
- dataTypes.go
- main.go

One bash file:
- format.sh


# External Libraries Used
[goreturns](https://github.com/sqs/goreturns) - Generator uses goreturns to remove unnecessary generated imports

[Substance](https://github.com/ahmedalhulaibi/substance) - This library is used to introspect on the database information schema and generate the graphql-go code.

[grahpql-go](https://github.com/graphql-go/graphql) - The generated code is using this implementation of GraphQL in Go.

[GORM](https://github.com/jinzhu/gorm) - The generated code is using GORM.

