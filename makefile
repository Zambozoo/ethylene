godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	cd src ; $$GOPATH/bin/godoc -http=:6060 -goroot . ; cd ..

run:
	cd src ; go run ./... $(ARGS) ; cd ..

test:
	cd src ; go test ./... -count=1 $(ARGS) ; cd ..