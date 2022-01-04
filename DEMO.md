# Observability Demo

Model for processing telemetry data

## Logs
Generate randomized logs at variable frequencies and with various decorators (account, service, etc)
> go run cmd/logs/main.go

## Stats
Man-in-the-middle statsd proxy - Logical stats flow: Service -> statsd proxy -> dogstatsd -> Datadog
> go run cmd/stats/main.go

## Tee
Read from a single source and write to two destinations (archives + forwarding)
> go run cmd/tee/main.go

## Sources and Outputs
1. Start Zookeeper and Kafka, along with the console consumer
> ./demo/restartKafka.sh

1. Run Kafka Log Demo:
> go run cmd/logs/kafka /main.go
