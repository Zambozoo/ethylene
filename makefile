godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	$$GOPATH/bin/godoc -http=:6060 -goroot .

run:
	go run ./... $(ARGS)

test:
	go test ./... -count=1 $(ARGS)