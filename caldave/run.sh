#!/bin/bash

# Kill any running instance of your Go server
pkill -f 'go run main.go'

# Run your Go server
go run main.go &
