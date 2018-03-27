# substance

_This project is a work in progress. This means the master branch is unstable and broken builds may occur._ 

Substance provides an interface for implementations extract schema data from a data provider (mysql, postgres, etc.).

SubstanceGen provides an interface for implementations to use data produced by Substance to generate code.

## Status

| Build                             | Code coverage                           | Report card                           | Codebeat                           |
| :-------------------------------: | :-------------------------------------: | :-----------------------------------: | :-----------------------------------: |
| [![Build][build-badge]][build-link] | [![Code coverage][cov-badge]][cov-link] | [![Report card][rc-badge]][rc-link]   | [![Codebeat][codebeat-badge]][codebeat-link]   |

[build-badge]: https://travis-ci.org/ahmedalhulaibi/substance.svg?branch=master "Travis-CI build status"
[build-link]: https://travis-ci.org/ahmedalhulaibi/substance "Travis-CI build status link"
[cov-badge]: https://codecov.io/gh/ahmedalhulaibi/substance/branch/master/graph/badge.svg "Code coverage status"
[cov-link]: https://codecov.io/gh/ahmedalhulaibi/substance "Code coverage status"
[rc-badge]: https://goreportcard.com/badge/github.com/ahmedalhulaibi/substance "Report card status"
[rc-link]: https://goreportcard.com/report/github.com/ahmedalhulaibi/substance "Report card status"
[codebeat-badge]: https://codebeat.co/badges/490b4031-5ae5-4fb5-bc46-b2c8802b944f
[codebeat-link]: https://codebeat.co/projects/github-com-ahmedalhulaibi-substance-master

## Substance Supported Data Providers

- MySQL/MariaDB
- Postgres
- JSON _planned_

## SubstanceGen Code Generators

- GraphQL-Go Server 
- GoStructs 
- GORM CRUD functions_WIP_
