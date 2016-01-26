all:
	go build -o kbartcheck cmd/kbartcheck/main.go
	go build -o holdingscov cmd/holdingscov/main.go

clean:
	rm -f ./kbartcheck
	rm -f ./holdingscov

test:
	go test -v ./...

bench:
	go test -bench=.

imports:
	goimports -w .
