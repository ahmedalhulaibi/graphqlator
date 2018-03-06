# Contributing

This is my first open source project. For now the contributing guidelines will be loose and we'll develop it as we go.


## Contributing to the Contributing Guide

If there is a good practice that anyone wants to add to the contributing guide for the sustainability of the project then by all means open an issue and we will have an open discussion.


# Pull Requests

Make sure to have tests up to date. Please include any environmental setup instructions.

Currently I am the only reviewer. I will do my best to review within 1-2 weeks.

## Current Environment Setup

The .travis.yml file contains the environment setup in terms of create the database schemas currently used.

# Issues

Open to any assistance in triaging and labelling issues. Hopefully we do not encounter too many problems/incidents.

## IMPORTANT

Graphqlator tries to purely remain a CLI program. All the database introspection and code generation is done in the [Substance](https://github.com/ahmedalhulaibi/substance) package in [vendor](https://github.com/ahmedalhulaibi/graphqlator/tree/master/vendor/github.com/ahmedalhulaibi/substance). 

If an enhancement/issue is related to the generated code itself, it is likely related to [Substance](https://github.com/ahmedalhulaibi/substance). For the time being we will still accept issues and enhancement requests for [Substance](https://github.com/ahmedalhulaibi/substance) here. May need a way to duplicate the issues automatically and label them on the [Substance](https://github.com/ahmedalhulaibi/substance) project as coming from graphqlator.

## Bugs or Unexpected Behaviour

Please write the steps to reproduce an issue. Make sure to include which database you are using in the title of the 

## Enhancements

No restrictions as of this writing on feature requests. When requesting an enhancement some use cases would be nice.

