gotable: *.go
	go vet
	golint
	go build

clean:
	rm -rf *.out

test:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out
	    
