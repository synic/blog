#!/bin/sh

HOST=${BLOG_HOST:=http://localhost:3000}
USERS=${NUM_USERS:=100}
COUNT=${REQUEST_COUNT:=150000}

if ! command -v cassowary &> /dev/null
then
    go install github.com/rogerwelin/cassowary/cmd/cassowary@v0.16.0
fi

cassowary run -u $BLOG_HOST \
    -c $USERS -n $COUNT -f scripts/.loadtest-urls.txt
