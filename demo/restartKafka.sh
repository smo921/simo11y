#!/usr/bin/env bash

ZOOKEEPER=/usr/local/opt/kafka/bin/zookeeper-server-start
KAFKA=/usr/local/opt/kafka/bin/kafka-server-start
BROKER="localhost:9092"
TOPIC="demo_topic"

if ! [[ -x ${ZOOKEEPER} || -x ${KAFKA} ]]; then
  echo "Zookeeper (${ZOOKEEPER} or Kafka (${KAFKA}) not found on system or not executable."
  exit 100
fi

# Cleanup previous Zookeeper / Kafka data
rm -rf /usr/local/var/lib/*

# Start Zookeeper locally:
${ZOOKEEPER} /usr/local/etc/kafka/zookeeper.properties > demo.log.zookeeper 2>&1 &
ZOOKEEPER_PID=$!
sleep 2
echo "Zookeeper started . . ."

# Start Kafka Broker:
${KAFKA} /usr/local/etc/kafka/server.properties > demo.log.kafka 2>&1 &
KAFKA_PID=$!
sleep 5 # NOTE: may need to increase sleep time if topic admin commands fail intermittently.
echo "Kafka started . . ."

trap "{ echo 'sigquit received' ; kill -9 $KAFKA_PID ; kill -9 $ZOOKEEPER_PID ; echo 'Zookeeper and Kafka terminated.'; rm -f demo.log.* ; }" QUIT

# Delete and create `demo_topic`:
kafka-topics --bootstrap-server ${BROKER} --topic ${TOPIC} --delete > /dev/null 2>&1
kafka-topics --bootstrap-server ${BROKER} --topic ${TOPIC} --create --partitions 3 --replication-factor 1 > /dev/null 2>&1
echo "Topic created: ${TOPIC}"

# Start Console consumer:
echo "Console consumer reading from ${TOPIC}"
kafka-console-consumer -brokers ${BROKER} -topic ${TOPIC} -offset oldest

kill -QUIT $$
