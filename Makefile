gotable: *.go
	go vet
	go build

clean:
	rm -rf *.out

lint:
	golint

test:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out
	    
