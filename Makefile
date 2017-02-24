gotable: *.go
	go clean
	go get -t ./...
	go vet
	go build

clean:
	go clean
	rm -rf *.out *.csv *.html *.txt *.pdf

lint:
	golint

test:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out

update:
	cp smoke_test.txt smoke_test.csv smoke_test.html testdata/
