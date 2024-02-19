#!/bin/sh
#
# A simple script to capture changes to any tailwind and templ files
# and run the server

pushd web/static
npx tailwindcss -i src/css/input.css -o ../dist/css/style.css

popd

pushd app/views
templ generate

popd

go run cmd/main.go
