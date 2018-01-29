# Graphqlator CLI
__*This project is a WIP. The end goal is to generate source for a graphql server*__

Graphqlator takes your existing database schema and generates code for a GraphQL-Go server. Type 'graphqlator help' to see usage.

A command line tool that generates a GraphQL type schema from a database table schema.

Supported Data Stores:

- mysql
- mariadb
- postgres

Install:

Download prebuilt binaries from the [releases page](https://github.com/ahmedalhulaibi/go-graphqlator-cli/releases)

Usage:
```
  graphqlator [flags]
  graphqlator [command]
```
Available Commands:
```
  describe    Describe database or table
  generate    Generate GraphQL type schema from database table.
  help        Help about any command
  version     Print the version number of Graphqlator
```
Flags:
  -h, --help   help for graphqlator

Use "graphqlator [command] --help" for more information about a command.

Example Usage:
```
graphqlator describe mysql "username:password@tcp(localhost:3306)/schema" OrdersTable PersonsTable AntiOrdersTable

graphqlator generate mysql "username:password@tcp(localhost:3306)/schema" OrdersTable PersonsTable AntiOrdersTable

graphqlator describe postgres "postgres://username:password@localhost:5432/postgres" orders persons antiorders

graphqlator generate postgres "postgres://username:password@localhost:5432/postgres" orders persons antiorders
```

# External Libraries
[Substance](https://github.com/ahmedalhulaibi/substance) - This library is used to introspect on the database information schema and generate the graphql-go code.

[grahpql-go](https://github.com/graphql-go/graphql) - The generated code is using this implementation of GraphQL in Go.

[GORM](https://github.com/jinzhu/gorm) - The generated code is using GORM.
