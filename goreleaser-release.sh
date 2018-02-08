#!/usr/bin/env bash
#check if HEAD is tagged
git describe --exact-match HEAD | grep -q fatal
grepRes=$?;
#if tagged then release to github
if [[ ! (( $grepRes )) ]]; 
then
    goreleaser --rm-dist
fi