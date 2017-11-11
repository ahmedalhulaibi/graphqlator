#!/usr/bin/env bash
go build
echo "TESTING ROOT AND HELP"
./graphqlator
./graphqlator help 
echo "TESTING DESCRIBES"
./graphqlator describe mysql "ahmed:password@tcp(localhost:3306)/awt_employee"
./graphqlator describe mysql "ahmed:password@tcp(localhost:3306)/awt_employee" Orders Persons
echo "TESTING GENERATES"
./graphqlator generate mysql "ahmed:password@tcp(localhost:3306)/awt_employee" Orders Persons