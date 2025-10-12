# 04textvar/

```console
go run 04textvar/main.go --help
Usage of textvar:
      --log-level    ENV: X_LOG_LEVEL   log level  {DEBUG, INFO, WARN, ERROR} (default INFO)
      --log-level2   ENV: X_LOG_LEVEL2  log level  {DEBUG, INFO, WARN, ERROR} (default INFO)
pflag: help requested
exit status 2

---

go run 04textvar/main.go --log-level DEBUG --log-level2 warn
parsed: &main.Options{LogLevel:-4, LogLevel2:(*slog.Level)(0x1400000e1d0)}
----------------------------------------
json:
{
  "LogLevel": "DEBUG",
  "LogLevel2": "WARN"
}
```
