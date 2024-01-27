#!/bin/bash -xe

# Fetch mastodon entries
go run src/main.go

# Build blog
cd blog
ls -R content
hugo
ls -R public