GO = go

pkg = sny.no/tools/edit
cmd = $(pkg)/cmd

edit = $(wildcard *.go)

.PHONY: all test fmt lint clean install

all: editd E B

test: lint
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...

clean:
	$(RM) editd E B

install:
	$(GO) install ./...

uninstall:
	$(RM) $(GOBIN)/editd $(GOBIN)/E $(GOBIN)/B

editd: $(wildcard cmd/editd/*.go) $(edit)
	$(GO) build $(cmd)/editd

E: $(wildcard cmd/E/*.go) $(edit) editd
	$(GO) build $(cmd)/E

B: $(wildcard cmd/B/*.go) $(edit) editd
	$(GO) build $(cmd)/B
