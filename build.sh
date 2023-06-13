#!/bin/bash

# Build the application
go build -o rent-app cmd/web/*.go  && ./rent-app
