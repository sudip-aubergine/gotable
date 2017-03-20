SCSS_BIN := scss

gotable: *.go css
	go clean
	go get -t ./...
	go vet
	go build

clean:
	go clean
	rm -rf *.out *.csv *.html *.txt *.pdf *.css* .sass-cache

css:
	${SCSS_BIN} ./scss/report.scss ./report.css --style=compressed --sourcemap=none
	@echo "Current working directory:"
	pwd
	@echo "scss completed.  ls -l ./report.css:"
	ls -l ./report.css

lint:
	golint

test:
	go test -v -coverprofile=coverage.out
	go tool cover -html=coverage.out

update:
	cp smoke_test.txt smoke_test.csv smoke_test.html smoke_test.pdf testdata/

all: clean gotable test

deps: wkhtmltopdf

wkhtmltopdf:
	./pdfinstall.sh
