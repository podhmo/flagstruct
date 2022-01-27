GO ?= go

00:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --verbose 

01:
	@diff -u 00*/main.go $@*/main.go > $(shell echo $@*/)/00$@.diff || :
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --verbose

	
