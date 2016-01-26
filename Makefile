all:
	go build -o kbartcheck cmd/kbartcheck/main.go

clean:
	rm -f ./kbartcheck
