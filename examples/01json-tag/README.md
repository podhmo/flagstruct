# 01json-tag/
```console
go run 01json-tag/main.go --help || :
Usage of hello:
      --another-name string   ENV: X_ANOTHER_NAME	-
      --name string           ENV: X_NAME	- (default "foo")
  -v, --verbose               ENV: X_VERBOSE	-
pflag: help requested
exit status 2


$ go run 01json-tag/main.go --verbose
parsed: &main.Options{Name:"foo", Verbose:true, Ignored:false, AnotherName:""}
```
