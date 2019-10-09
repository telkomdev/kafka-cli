## Kafka CLI

Kafka cli client

### Install

<b>Build from source</b>

Require `Go` version `1.10` or higher
```shell
$ git clone https://github.com/musobarlab/kafka-cli.git

$ make build

$ kafka-cli --version
```

<b>Mac OS</b>
```shell
$ brew tap wuriyanto48/tool

$ brew install kafka-cli

$ kafka-cli --version
```

<b>Linux</b>
```shell
$ wget https://github.com/musobarlab/kafka-cli/releases/download/v0.0.0/kafka-cli-v0.0.0.linux-amd64.tar.gz

$ tar -zxvf kafka-cli-v0.0.0.linux-amd64.tar.gz

$ kafka-cli --version
```

<b>Windows</b>

Download latest version https://github.com/musobarlab/kafka-cli/releases

### Usage

<b>Publish message to Kafka broker and topic</b>
```shell
$ kafka-cli pub -broker localhost:9092 -topic wurys -m "hahahaha" -V
```

<b>JSON</b>
```shell
$ kafka-cli pub -broker localhost:9092 -topic wurys -m "{"hello":"hello", "world":"world"}" -V
```

<b>or multiple broker</b>
```shell
$ kafka-cli pub -broker localhost:9092,localhost:9093,localhost:9094 -topic wurys -m "hahahaha" -V
```

<b>Subscribe to Kafka broker and topic</b>
```shell
$ kafka-cli sub -broker localhost:9092 -topic wurys
```

<b>or multiple broker</b>
```shell
$ kafka-cli sub -broker localhost:9092,localhost:9093,localhost:9094 -topic wurys
```