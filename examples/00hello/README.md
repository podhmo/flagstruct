# 00hello/
```console
go run 00hello/main.go --help || :
Usage of hello:
      --name string   ENV: X_NAME	name of greeting (default "foo")
  -v, --verbose       ENV: X_VERBOSE	-
pflag: help requested
exit status 2


$ go run 00hello/main.go --verbose 
parsed: &main.Options{Name:"foo", Verbose:true}
```
