# dalelane-kafka-datasource

A Grafana data source plugin for Apache Kafka - for visualising the contents of Kafka messages

## Configuration

| **Option**              | **Default** | **Explanation**                                            |
| ----------------------- | ----------- | ---------------------------------------------------------- |
| **Bootstrap servers**   |             | address to connect to (in the form of `host:port`)         |
| **Client id**           | `grafana`   | how Grafana should identify itself to Kafka                |
| **Consumer Group id**   | `grafana`   | name of the consumer group Grafana should join             |
| **Authentication type** | `none`      | SASL mechanism to use (None if no authentication required) |
| **username**            |             | SASL username (ignored if Authentication type is None)     |
| **password**            |             | SASL password (ignored if Authentication type is None)     |
| **Use TLS**             | `false`     | true for security protocol SASL_SSL                        |

## Screenshots

### Configuring a new data source

![screenshot](https://github.com/dalelane/grafana-kafka-datasource/raw/master/src/img/config.png)

### Dashboard plotting Kafka events as timeseries data

![screenshot](https://github.com/dalelane/grafana-kafka-datasource/raw/master/src/img/dashboard-1.png)

## Limitations / TODOs

- Only JSON data is currently supported (other formats, such as Apache Avro, to be added)
- Only unauthorised, or SASL/SCRAM auth types are currently supported (other auth mechanisms, such as SASL/PLAIN and mTLS, to be added)
