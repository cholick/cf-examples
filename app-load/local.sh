#!/usr/bin/env bash

export VCAP_APPLICATION='{
  "application_uris": [
   "localhost"
  ]
 }
'
export PORT=3000
export REQUEST_PORT=3000
export DELAY=250ms

go run main.go
