### Kafka CLI

Kafka cli client

#### Install

From source
```shell
$ go get github.com/musobarlab/kafka-cli

$ go install github.com/musobarlab/kafka-cli/cmd

$ kafka-cli --version
```

#### Usage

Publish message to Kafka broker and topic
```shell
$ kafka-cli pub -broker localhost:9092 -topic wurys -m "hahahaha" -V
```

or multiple broker
```shell
$ kafka-cli pub -broker localhost:9092,localhost:9093,localhost:9094 -topic wurys -m "hahahaha" -V
```

Subscribe to Kafka broker and topic
```shell
$ kafka-cli sub -broker localhost:9092 -topic wurys
```

or multiple broker
```shell
$ kafka-cli sub -broker localhost:9092,localhost:9093,localhost:9094 -topic wurys
```