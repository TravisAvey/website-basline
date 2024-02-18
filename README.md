# website-baseline
Website Baseline


### Tailwindcss notes:
npx tailwindcss -o dist/css/style.css --minify

### Postgres Podman notes
`podman run --name psql -p 5432:5432 -e POSTGRES_USER=psqluser -e POSTGRES_PASSWORD=password -e POSTGRES_DB=web -v /website-baseline/config/init.sql:/docker-entrypoint-initdb.d/init.sql -d postgres:16`

to look inside and take a peak:
`
podman exec -it psql psql -U psqluser -d web
`
