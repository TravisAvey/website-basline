#!/bin/sh
#
# A simple script to capture changes to any tailwind and templ files
# and run the server

pushd web/static
npm run build
npm run alpine
npm run scripts
npm run editor
popd

go run cmd/main.go
