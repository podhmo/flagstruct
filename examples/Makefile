SHELL := bash
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
	
04:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --ip 127.0.0.1
	
05:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --value 20
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --value foo || : # ng
	@echo ---
	$(GO) run $(shell echo $@*/)main.go # ng

06:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --name foo
	@echo ---
	$(GO) run $(shell echo $@*/)main.go # ng

07:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --name foo
	@echo ---
	$(GO) run $(shell echo $@*/)main.go # ng

08:
	$(GO) run $(shell echo $@*/)main.go --help || :
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --a.name aaa --b.verbose
	@echo ---
	$(GO) run $(shell echo $@*/)main.go --a.name aaa --b.verbose --debug

doc:
	$(info 00)
	printf '# $(shell echo 00*/)\n' > $(shell echo 00*/)README.md
	printf '```console\n$ ' >> $(shell echo 00*/)README.md
	$(MAKE) 00 |& grep -vP '^make' | perl -p -e 's/---\n/\n\n\$$ /' >> $(shell echo 00*/)README.md
	printf '```\n' >> $(shell echo 00*/)README.md
	$(info 01)
	printf '# $(shell echo 01*/)\n' > $(shell echo 01*/)README.md
	printf '```console\n$ ' >> $(shell echo 01*/)README.md
	$(MAKE) 01 |& grep -vP '^make' | perl -p -e 's/---\n/\n\n\$$ /' >> $(shell echo 01*/)README.md
	printf '```\n' >> $(shell echo 01*/)README.md
	$(info 02)
	printf '# $(shell echo 02*/)\n' > $(shell echo 02*/)README.md
	printf '```console\n$ ' >> $(shell echo 02*/)README.md
	$(MAKE) 02 |& grep -vP '^make' | perl -p -e 's/---\n/\n\n\$$ /' >> $(shell echo 02*/)README.md
	printf '```\n' >> $(shell echo 02*/)README.md
	$(info 03)
	printf '# $(shell echo 03*/)\n' > $(shell echo 03*/)README.md
	printf '```console\n$ ' >> $(shell echo 03*/)README.md
	$(MAKE) 03 |& grep -vP '^make' | perl -p -e 's/---\n/\n\n\$$ /' >> $(shell echo 03*/)README.md
	printf '```\n' >> $(shell echo 03*/)README.md
	$(info 04)
	printf '# $(shell echo 04*/)\n' > $(shell echo 04*/)README.md
	printf '```console\n$ ' >> $(shell echo 04*/)README.md
	$(MAKE) 04 |& grep -vP '^make' | perl -p -e 's/---\n/\n\n\$$ /' >> $(shell echo 04*/)README.md
	printf '```\n' >> $(shell echo 04*/)README.md
	$(info 08)
	printf '# $(shell echo 08*/)\n' > $(shell echo 08*/)README.md
	printf '```console\n$ ' >> $(shell echo 08*/)README.md
	$(MAKE) 08 |& grep -vP '^make' | perl -p -e 's/---\n/\n\n\$$ /' >> $(shell echo 08*/)README.md
	printf '```\n' >> $(shell echo 08*/)README.md