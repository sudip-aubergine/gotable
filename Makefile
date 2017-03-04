SCSS_BIN := scss

gotable: *.go
	go clean
	go get -t ./...
	go vet
	go build

clean:
	go clean
	rm -rf *.out *.csv *.html *.txt *.pdf *.css* .sass-cache

css:
	${SCSS_BIN} ./scss/report.scss ./report.css --style=expanded --sourcemap=none

lint:
	golint

test:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

update:
	cp smoke_test.txt smoke_test.csv smoke_test.html testdata/
