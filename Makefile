CLI = ./bin/uhppote-cli
SIMULATOR = ./bin/uhppote-simulator
DEBUG = --debug
LOCAL = 192.168.1.100:51234
CARD = 6154410
SERIALNO = 423187757
DOOR = 3
VERSION = 0.03.0
DIST ?= development

all: test      \
	 benchmark \
     coverage

format: 
	go fmt uhppote...

release: format
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/openapi
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppote-cli             uhppote-cli
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppote-cli            uhppote-cli
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppote-cli.exe       uhppote-cli
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppoted-rest           uhppoted-rest
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppoted-rest          uhppoted-rest
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppoted-rest.exe     uhppoted-rest
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppoted-mqtt           uhppoted-mqtt
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppoted-mqtt          uhppoted-mqtt
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppoted-mqtt.exe     uhppoted-mqtt
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppote-simulator       uhppote-simulator
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppote-simulator      uhppote-simulator
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppote-simulator.exe uhppote-simulator
	cp -r install/openapi/* dist/$(DIST)/openapi/

build: format
	go install uhppote-cli
	go install uhppote-simulator
	go install uhppoted-rest
	go install uhppoted-mqtt

test: build
	go clean -testcache
	go test uhppote...

integration-tests: build
	go clean -testcache
	go test -count=1 src/integration-tests/*.go

benchmark: build
	go test src/uhppote/encoding/bcd/*.go -bench .

coverage: build
	go clean -testcache
	go test -cover src/uhppote/*.go
	go test -cover src/uhppote/encoding/bcd/*.go
	go test -cover src/uhppote/encoding/UTO311-L0x/*.go
	go test -cover src/uhppoted-rest/encoding/plist/*.go

clean:
	go clean
	rm -rf bin

usage: build
	$(CLI)

debug: build
	go clean -testcache
	go test -v src/uhppote/encoding/conf/*.go
	go test -v src/uhppoted-rest/config/*.go
	go test -v src/uhppoted-mqtt/config/*.go

help: build
	$(CLI)       help
	$(CLI)       help get-devices
	$(SIMULATOR) help
	$(SIMULATOR) help new-device

version: build
	$(CLI)       version
	$(SIMULATOR) version

run: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-devices
#	$(CLI) --bind $(LOCAL) $(DEBUG) set-address    $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'
	$(CLI) --bind $(LOCAL) $(DEBUG) get-cards      $(SERIALNO)
	$(CLI) --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)
	$(CLI) --bind $(LOCAL) $(DEBUG) set-time       $(SERIALNO)
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke         $(SERIALNO) $(CARD)

get-devices: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-devices

set-address: build
	$(CLI) -bind $(LOCAL) $(DEBUG) set-address $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'

get-status: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-status $(SERIALNO)

get-time: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-time $(SERIALNO)

set-time: build
	# $(CLI) -debug set-time 423187757 '2019-01-08 12:34:56'
	$(CLI) --bind $(LOCAL) $(DEBUG) set-time $(SERIALNO)

get-door-delay: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)

set-door-delay: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-door-delay $(SERIALNO) $(DOOR) 5

get-listener: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-listener $(SERIALNO)

set-listener: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-listener $(SERIALNO) 192.168.1.100:40000

get-cards: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-cards $(SERIALNO)

get-card: build
	$(CLI) $(DEBUG) get-card $(SERIALNO) $(CARD)

grant: build
	$(CLI) $(DEBUG) grant $(SERIALNO) $(CARD) 2019-01-01 2019-12-31 1,2,3,4

revoke: build
	$(CLI) $(DEBUG) revoke $(SERIALNO) $(CARD)

revoke-all: build
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke-all $(SERIALNO)

load-acl: build
	$(CLI) --config .UTO311-L04 $(DEBUG) load-acl debug.tsv

get-events: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-events $(SERIALNO)

get-event-index: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-event-index $(SERIALNO)

set-event-index: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-event-index $(SERIALNO) 23

open: build
	$(CLI) --bind $(LOCAL) $(DEBUG) open $(SERIALNO) 1

listen: build
	$(CLI) --bind 192.168.1.100:40000 $(DEBUG) listen 

simulator: build
	./bin/uhppote-simulator --debug --devices "./runtime/simulation/devices"

simulator-device: build
	./bin/uhppote-simulator --debug --devices "runtime/simulation/devices" new-device 666

uhppoted-rest: build
	./bin/uhppoted-rest --console

uhppoted-rest-daemonize: build
	sudo ./bin/uhppoted-rest daemonize

uhppoted-rest-undaemonize: build
	sudo ./bin/uhppoted-rest undaemonize

uhppoted-rest-version: build
	./bin/uhppoted-rest version

uhppoted-rest-help: build
	./bin/uhppoted-rest help
	./bin/uhppoted-rest help commands
	./bin/uhppoted-rest help version
	./bin/uhppoted-rest help help

uhppoted-rest-linux: build
	mkdir -p ./dist/development/linux
	env GOOS=linux GOARCH=amd64 go build -o dist/development/linux/uhppoted-rest uhppoted-rest

uhppoted-rest-windows: build
	mkdir -p ./dist/development/windows
	env GOOS=windows GOARCH=amd64 go build -o dist/development/windows/uhppoted-rest.exe uhppoted-rest

uhppoted-rest-docker: build
	env GOOS=linux GOARCH=amd64 go build -o docker/linux/uhppote-simulator uhppote-simulator
	env GOOS=linux GOARCH=amd64 go build -o docker/linux/uhppoted-rest     uhppoted-rest
	docker build -f ./docker/Dockerfile.uhppoted -t uhppoted . 
	docker run --detach --publish 8080:8080 --rm uhppoted

uhppoted-mqtt: build
	./bin/uhppoted-mqtt --console

uhppoted-mqtt-help: build
	./bin/uhppoted-mqtt help
	./bin/uhppoted-mqtt help commands
	./bin/uhppoted-mqtt help version
	./bin/uhppoted-mqtt help help

uhppoted-mqtt-version: build
	./bin/uhppoted-mqtt version

swagger: 
	docker run --detach --publish 80:8080 --rm swaggerapi/swagger-editor 

hivemq:
	docker run --detach --publish 8080:8080 --publish 1883:1883 --rm hivemq/hivemq4
