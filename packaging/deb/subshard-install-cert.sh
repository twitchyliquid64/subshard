#!/bin/bash

certutil -d "sql:$HOME/.pki/nssdb" -A -n subshard -i "$1" -t C
