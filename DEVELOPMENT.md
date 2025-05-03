# dalelane-kafka-datasource

## Dependencies

- Docker (with Docker Compose)
- Go (I used 1.24.2)
- Node (I used 20.10.0)

## Environment

A Docker Compose environment is included to enable local development.

It includes:

- `kafka` : a one-broker Kafka cluster
- `connect` : a Kafka Connect instance to generate events in Kafka
- `grafana` : instance of Grafana (with the plugin loaded)
- `renderer` : Grafana image renderer

Ports used:

- (kafka) `localhost:9092` - bootstrap address for the Kafka cluster
- (grafana) http://localhost:3000 - web address for the Grafana server

## Scripts

To build from source:

```sh
./scripts/build.sh
```

To run the development build in a Docker Compose environment:

```sh
./scripts/run.sh
```

To delete everything:

```sh
./scripts/clean.sh
```
