NAME 		:= haki
BIN        	:= $(NAME)
REPO        := github.com/rjansen/$(NAME)
BUILD       := $(shell git rev-parse --short HEAD)
VERSION     := $(shell git describe --tags --always)
MAKEFILE    := $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    := $(shell cd $(dir $(MAKEFILE)); pwd)
ALLPKGS    	:= $(shell go list ./... | grep -v /vendor/)
PKGS       	:= $(shell echo $(ALLPKGS) | grep -v /itest)
IPKGS      	:= $(shell echo $(ALLPKGS) | grep /itest)

# Test and Benchmark Parameters
TEST_PKGS ?=
TESTS ?= .
COVERAGE_FILE := $(NAME).coverage
COVERAGE_HTML := $(NAME).coverage.html
PKG_COVERAGE := $(NAME).pkg.coverage

ENV ?= local

.PHONY: default
default: version

.PHONY: version
version:
	@echo "Version: $(REPO)@$(VERSION)-$(BUILD)"
	@echo "AllPkgs: $(ALLPKGS)"
	@echo "Pkgs: $(PKGS)"
	@echo "IPkgs: $(IPKGS)"

.PHONY: install
install: install.tools install.deps
	@echo "$(REPO) installed successfully"

.PHONY: install.deps
install.deps: workspace
	go get -u github.com/Masterminds/glide
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: workspace
workspace: install.tools
	gvm use go1.10.1

.PHONY: install.tools
install.tools:
	#if [[ ! which gvm ]]; then \
	   # bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
	   # gmv install go1.10.1; \
	#fi
	#Create abstraction to install into linux distro too
	#if [[ ! $$(which direnv) ]]; then \
	#	brew install direnv; \
	#fi
	#Create abstraction to install into linux distro too
	#if [[ ! $$(which direnv) ]]; then \
	#	curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sudo sh; \
	#fi

.PHONY: setup
setup: workspace
	@echo "$(REPO) setup successfully"
	glide sync

.PHONY: all
all: build test bench coverage

.PHONY: build
build: test
	@echo "Building: $(REPO)@$(VERSION)-$(BUILD)"
	go build $(REPO)

.PHONY: build
buildi.anyway:
	@echo "Try to Building: $(REPO)@$(VERSION)-$(BUILD)"
	go build $(REPO)

.PHONY: deploy
deploy:
	@echo "Deploying: $(REPO)@$(VERSION)-$(BUILD)"

.PHONY: default
default: build

local:
	@echo "Set enviroment to local"
	$(eval ENV = local)

.PHONY: dev
dev:
	@echo "Set enviroment to dev"
	$(eval ENV = dev)

.PHONY: prod
prod:
	@echo "Set enviroment to prod"
	$(eval ENV = prod)

.PHONY: check_env
check_env:
	@if [ "$(ENV)" == "" ]; then \
	    echo "Env is blank: $(ENV)"; \
	    exit 540; \
	fi

.PHONY: run
run: build
	@echo "Running: $(REPO)@$(VERSION)-$(BUILD)"

.PHONY: debug
debug: build
	@echo "Debugging: $(REPO)@$(VERSION)-$(BUILD)"
	#gdb ./$(NAME)

.PHONY: clean
clean:
	-rm $(NAME)*coverage*
	-rm *.test
	-rm *.pprof

.PHONY: reset
reset: clean
	-cd vendor; rm -r */

.PHONY: test
test:
	@#go test -v -race $(ALLPKGS)
	@echo "Testing: $(REPO)@$(VERSION)-$(BUILD)"
	@if [ "$(TEST_PKGS)" = "" ]; then \
	    echo "Unit Test All Pkgs with Selected Tests=$(TESTS)"; \
		go test -v -race -run $(TESTS) $(PKGS) || exit 501; \
	else \
	    echo "Unit Test Selected Pkgs=$(TEST_PKGS) with Selected Tests=$(TESTS)"; \
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -race -run $(TESTS) $(REPO)/$$tstpkg || exit 501; \
		done; \
	fi

.PHONY: itest
itest:
	@echo "Integreation Testing: $(REPO)@$(VERSION)-$(BUILD)"
	@if [ "$(TEST_PKGS)" = "" ]; then \
	    echo "Integration Test All Pkgs";\
		go test -v -race $(IPKGS) || exit 501;\
	else \
	    echo "Integration Unit Test Selected Pkgs=$(TEST_PKGS)";\
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -race $(REPO)/$$tstpkg/itest || exit 501;\
		done; \
	fi

.PHONY: bench
bench:
	@#go test -bench=. -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $(IPKGS)
	@#go test -bench=. -run="^$$" -benchmem $(PKGS)
	@echo "Benchmarking: $(REPO)@$(VERSION)-$(BUILD)"
	@if [ "$(TEST_PKGS)" = "" ]; then \
	    echo "Benchmark all Pkgs" ;\
	    for tstpkg in $(IPKGS); do \
		    go test -bench=. -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $$tstpkg || exit 501;\
		done; \
	else \
	    echo "Benchmark Selected Pkgs=$(TEST_PKGS)" ;\
	    for tstpkg in $(TEST_PKGS); do \
		    go test -bench=. -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $(REPO)/$$tstpkg/itest || exit 501;\
		done; \
	fi

.PHONY: coverage
coverage:
	@echo "Coveraging: $(REPO)@$(VERSION)-$(BUILD)"
	@echo 'mode: set' > $(COVERAGE_FILE)
	@touch $(PKG_COVERAGE)
	@touch $(COVERAGE_FILE)
	@if [ "$(TEST_PKGS)" == "" ]; then \
		for pkg in $(ALLPKGS); do \
			go test -v -coverprofile=$(PKG_COVERAGE) $$pkg || exit 501; \
			if (( `grep -c -v 'mode: set' $(PKG_COVERAGE)` > 0 )); then \
				grep -v 'mode: set' $(PKG_COVERAGE) >> $(COVERAGE_FILE); \
			fi; \
		done; \
	else \
	    echo "Testing with covegare the Pkgs=$(TEST_PKGS)" ;\
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -coverprofile=$(PKG_COVERAGE) $(REPO)/$$tstpkg || exit 501; \
			if (( `grep -c -v 'mode: set' $(PKG_COVERAGE)` > 0 )); then \
				grep -v 'mode: set' $(PKG_COVERAGE) >> $(COVERAGE_FILE); \
			fi; \
		done; \
	fi
	@echo "Generating report"
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	open $(COVERAGE_HTML) || google-chrome $(COVERAGE_HTML)

.PHONY: docker.build
docker.build:
	@echo "Building Docker: $(REPO)@$(VERSION)-$(BUILD)"
	@docker build -t $(NAME) .

#%.alias:
#	@apex alias $(ALIAS) $*
#	@echo '---- alias function $*:$(ALIAS) ----'


docker.%: docker.build
	@echo "Executing: $* on Docker: $(REPO)@$(VERSION)-$(BUILD)"
	@docker run --rm --name $(NAME) -w="/go/src/$(REPO)" -v `pwd`:/go/src/$(REPO) $(NAME) make $*
