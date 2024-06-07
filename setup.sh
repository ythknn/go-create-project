#!/bin/bash

# Build the create project
sudo go build -o go-create-project

# Move the binary to /usr/local/bin/
sudo mv go-create-project /usr/local/bin/
