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

02:
	$(GO) run $(shell echo $@*/)main.go --log-level DEBUG --log-level2 warn
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --log-level foo --log-level2 foo

03:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	DB_URI=sqlite:///data.db $(GO) run $(shell echo $@*/)main.go --db.debug --another-db.uri sqlite:///:memory:
	
