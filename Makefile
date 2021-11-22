.PHONY : test build clean format

build:
	go build github.com/telkomdev/kafka-cli/cmd/kafka-cli

test:
	go test ./...

format:
	find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" | xargs gofmt -s -d -w

clean:
	rm kafka-cli