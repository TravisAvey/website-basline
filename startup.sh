#!/bin/sh
#
# A simple script to capture changes to any tailwind and templ files
# and run the server

pushd web/static
npx tailwindcss -i src/css/input.css -o ../dist/css/style.css

popd

go run cmd/main.go
