GOVENDOR=$(shell echo "$(GOBIN)/govendor")

test: $(GOVENDOR) vet
	env $(shell cat .env.test) govendor test -v +local

vet: $(GOVENDOR)
	govendor vet +local

$(GOVENDOR):
	go get -v github.com/kardianos/govendor
