# 08shared-common-option/
```console
go run 08shared-common-option/main.go --help || :
Usage of /tmp/go-build3591807064/b001/exe/main:
      --a.name string   ENV: A_NAME	-
      --b.verbose       ENV: B_VERBOSE	-
      --debug           ENV: DEBUG	-
pflag: help requested
exit status 2


$ go run 08shared-common-option/main.go --a.name aaa --b.verbose
{
  "debug": false,
  "a": {
    "debug": false,
    "name": "aaa"
  },
  "b": {
    "debug": false,
    "verbose": true
  }
}


$ go run 08shared-common-option/main.go --a.name aaa --b.verbose --debug
{
  "debug": true,
  "a": {
    "debug": true,
    "name": "aaa"
  },
  "b": {
    "debug": true,
    "verbose": true
  }
}
```
