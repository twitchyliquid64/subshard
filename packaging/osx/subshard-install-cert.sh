#!/bin/bash

security add-trusted-cert -p ssl -r trustRoot -k ~/Library/Keychains/login.keychain "$1"
