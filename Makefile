SCSS_BIN := scss

gotable: css defaults *.go
	go clean
	go get -t -v ./...
	go vet
	go build

clean:
	go clean
	rm -rf *.out *.csv *.html *.txt *.pdf *.css* .sass-cache
	rm -f defaults.go

defaults:
	./defaults.sh

css:
	${SCSS_BIN} ./scss/gotable.scss ./gotable.css --style=compressed --sourcemap=none
	@echo "Current working directory:"
	pwd
	@echo "scss completed.  ls -l ./gotable.css:"
	ls -l ./gotable.css

lint:
	golint

test:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

update:
	cp smoke_test.txt smoke_test.csv smoke_test.html smoke_test.pdf smoke_test_custom_template.html testdata/

all: clean gotable test

deps: wkhtmltopdf

wkhtmltopdf:
	./pdfinstall.sh
