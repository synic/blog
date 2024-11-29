#!/bin/sh

HOST=${BLOG_HOST:=http://localhost:3000}
USERS=${NUM_USERS:=100}
COUNT=${REQUEST_COUNT:=150000}

go run github.com/rogerwelin/cassowary/cmd/cassowary@latest run \
    -u $BLOG_HOST \
    -c $USERS \
    -n $COUNT \
    -f scripts/.loadtest-urls.txt \
    -H "Accept-Encoding: identity"
