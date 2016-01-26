all:
	go build -o kbartcheck cmd/kbartcheck/main.go
	go build -o holdingcov cmd/holdingcov/main.go

clean:
	rm -f ./kbartcheck
	rm -f ./holdingcov

test:
	go test -v ./...

bench:
	go test -bench=.
