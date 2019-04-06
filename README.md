so-postgres-compose
===================

This repository contains code and config files for my answer to the question [Cannot connect to db host within docker containers from api-service to db-service in order to do migration using goose in golang][so:q]


## Prequisites

* docker/[Docker Desktop][docker:desktop]
* [docker-compose][docker:compose]

## Usage

``` sh
git clone https://github.com/mwmahlberg/so-postgres-compose.git
cd so-postgres-compose
docker-compose build
docker-compose up -d
```

Then, you can open <http://localhost:8081>

[docker:desktop]: https://www.docker.com/products/docker-desktop
[docker:compose]: https://docs.docker.com/compose/
[so:q]: https://stackoverflow.com/questions/55547207/cannot-connect-to-db-host-within-docker-containers-from-api-service-to-db-servic