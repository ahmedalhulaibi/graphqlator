#!/bin/sh

newtag=$(grep -Po 'v[0-9]+\.[0-9]+\.[0-9]+\-?([a-z]*)?' cmd/versionnumber.go)
git tag $newtag