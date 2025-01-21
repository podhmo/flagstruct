# 02enum/

enum example

```console
go run 02enum/main.go --log-level DEBUG --log-level2 warn
parsed: &main.Options{Name:"foo", LogLevel:"DEBUG", LogLevel2:(*main.LogLevel)(0xc000098600)}


$ go run 02enum/main.go --log-level foo --log-level2 foo
invalid argument "foo" for "--log-level" flag: FOO is an invalid value for main.LogLevel
Usage of hello:
      --log-level LogLevel    ENV: X_LOG_LEVEL	log level {DEBUG, INFO, WARN, ERROR} (default INFO)
      --log-level2 LogLevel   ENV: X_LOG_LEVEL2	log level {DEBUG, INFO, WARN, ERROR} (default INFO)
      --name string           ENV: X_NAME	name of greeting (default "foo")
invalid argument "foo" for "--log-level" flag: FOO is an invalid value for main.LogLevel
exit status 2
```

more simplified definition with TextVar is [here](./with-textvar)