
dependency:
	go get -u github.com/stretchr/testify/assert

test:
	go test -v -race

bench:
	go test -bench=. -race
