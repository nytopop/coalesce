#!/bin/sh

#reflex -s -r '.*\.\(go\|html\)' ./build.sh

# '\.go$'

# '.*\.\(go\|html\)'

reflex -r '\.go$' -s -- sh -c 'go fmt *.go; go build -o bin/coalesce && bin/coalesce'
