# Observability Demo

Model for processing telemetry data

## Logs
> go run cmd/logs/main.go

## Stats
> go run cmd/stats/main.go

## Tee
> go run cmd/tee/main.go

## Sources and Syncs

Cleanup previous Zookeeper / Kafka data
> rm -rf /usr/local/var/lib/*

Start Zookeeper locally:
> /usr/local/opt/kafka/bin/zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties

Start Kafka Broker:
> /usr/local/opt/kafka/bin/kafka-server-start /usr/local/etc/kafka/server.properties

Delete and create `demo_topic`:
> kafka-topics --bootstrap-server localhost:9092 --topic demo_topic --delete
> kafka-topics --bootstrap-server localhost:9092 --topic demo_topic --create --partitions 3 --replication-factor 1

Start Console consumer:
> kafka-console-consumer -brokers localhost:9092 -topic demo_topic -offset oldest

Run Kafka Log Demo:
> go run cmd/logs/kafka /main.go
