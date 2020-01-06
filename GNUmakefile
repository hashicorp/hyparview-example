# ======================================================================
# Entry Points

build: bin/hyparview-example
.PHONEY: build

run: clientCount = 10
run: clientNums = $(shell seq $(clientCount))
run: clients = $(foreach c,$(clientNums),data/$(c).log)
run: bin/hyparview-example
	echo $(clients)
.PHONEY: run

# ======================================================================
# Implementation

# build
sources = $(wildcard *.go) $(wildcard proto/*.go)
bin/hyparview-example: $(sources)
	go build
	mv hyparview-example $@

# certs & keys
domain = test

cert/$(domain)-agent-ca.pem:
	mkdir -p cert
	cd cert && consul tls ca create -domain=$(domain)

cert/$(domain)-server-%.pem: cert/$(domain)-agent-ca.pem
	cd cert && consul tls cert create -server -domain $(domain)
	mv cert/dc1-server-$(domain)-0.pem $@
	mv cert/dc1-server-$(domain)-0-key.pem $(basename $@)-key.pem

cert/$(domain)-client-%.pem: cert/$(domain)-agent-ca.pem
	cd cert && consul tls cert create -client -domain $(domain)
	mv cert/dc1-client-$(domain)-0.pem $@
	mv cert/dc1-client-$(domain)-0-key.pem $(basename $@)-key.pem

# run clients
data/%.log: ADDR = localhost:444$*
data/%.log: BOOTSTRAP = localhost:4440
data/%.log: SERVER_PEM = cert/$(domain)-server-$*.pem
data/%.log: SERVER_KEY = cert/$(domain)-server-$*-key.pem
data/%.log: CLIENT_PEM = cert/$(domain)-client-$*.pem
data/%.log: CLIENT_KEY = cert/$(domain)-client-$*-key.pem
export ADDR BOOTSTRAP SERVER_PEM SERVER_KEY CLIENT_PEM CLIENT_KEY
data/%.log: bin/hyparview-example data $(SERVER_PEM) $(CLIENT_PEM)
	bin/hyparview-example > $@ 2>&1 &!
