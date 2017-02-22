gotable: *.go
	go get -t ./...
	go vet
	go build

clean:
	rm -rf *.out *.csv *.html

lint:
	golint

test:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out
	    
