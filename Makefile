.DEFAULT_GOAL := all

.PHONY: all
all: fmt tidy

.PHONY: tidy
tidy:
	@go mod tidy -go=1.16 && go mod tidy -go=1.17

.PHONY: fmt
fmt:
	@gofmt -w -e -s -l .

.PHONY: commit
commit: fmt
	@git add . && git commit -m "ok"

.PHONY: push
push: commit
	@git push