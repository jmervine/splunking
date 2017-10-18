GOVENDOR=$(shell echo "$(GOBIN)/govendor")

test: $(GOVENDOR) vet
	govendor test +local

vet: $(GOVENDOR)
	govendor vet +local

$(GOVENDOR):
	go get -v github.com/kardianos/govendor
