# ======================================================================
# Entry Points

build: bin/hyparview-example
.PHONEY: build

root: cert/$(DOMAIN)-agent-ca.pem
.PHONEY: root

# ======================================================================
# Implementation

# build
sources = $(wildcard *.go)
bin/hyparview-example: $(sources) proto/hyparview.pb.go proto/gossip.pb.go
	go build
	mv hyparview-example $@

proto/%.pb.go: proto/%.proto
	protoc -I . --go_out=plugins=grpc:. $<

# certs & keys
cert/$(DOMAIN)-agent-ca.pem:
	mkdir -p cert
	cd cert && consul tls ca create -domain=$(DOMAIN)

cert/$(DOMAIN)-server-%.pem: cert/$(DOMAIN)-agent-ca.pem
	cd cert && consul tls cert create -server -domain $(DOMAIN)
	mv cert/dc1-server-$(DOMAIN)-0.pem $@
	mv cert/dc1-server-$(DOMAIN)-0-key.pem $(basename $@)-key.pem

cert/$(DOMAIN)-client-%.pem: cert/$(DOMAIN)-agent-ca.pem
	cd cert && consul tls cert create -client -domain $(DOMAIN)
	mv cert/dc1-client-$(DOMAIN)-0.pem $@
	mv cert/dc1-client-$(DOMAIN)-0-key.pem $(basename $@)-key.pem
