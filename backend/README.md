# Backend

TODO.

## Setting up the environment

### PostgreSQL

- https://hub.docker.com/\_/postgres
- https://pkg.go.dev/github.com/lib/pq

`docker run \
	--name some-postgres \
    -e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-p 5432:5432 \
	-d postgres`

### (Optional) PgAdmin

`docker run -d \
    --name some-phppgadmin \
    --link some-postgres:postgres \
    -e 'PGADMIN_DEFAULT_EMAIL=user@domain.com' \
    -e 'PGADMIN_DEFAULT_PASSWORD=supersecret' \
    -e disable_ssl=true \
    -p 8080:80 \
    dpage/pgadmin4`

Verify that it is running via `docker ps` and grab the IP address of the container.
`docker inspect some-postgres -f "{{json .NetworkSettings.Networks.bridge.IPAddress }}"`

The web interface will be hosted at `http://10.0.2.3:8080`.

You will now need to add a configuration that allows pgAdmin to connect to
the PostgreSQL container. Follow the steps below.

Servers > Register > Server...

- Name: any
  Connection
- Host: Host of the machine
- Username: postgres
- Password: postgres
  Press Save
