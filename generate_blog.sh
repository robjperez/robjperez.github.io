#!/bin/bash

# Fetch mastodon entries
go run src/main.go

# Build blog
cd blog && hugo