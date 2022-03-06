.PHONY: clean test coverage
.SILENT: ${.PHONY}

GO_TEST := go test ./... -race -p 1

all: clean
	${GO_TEST} -v

clean:
	rm -f *.out *.html
	go clean -cache

test: clean
	${GO_TEST}

coverage: clean
	${GO_TEST} -v -covermode atomic -coverprofile cover.out
