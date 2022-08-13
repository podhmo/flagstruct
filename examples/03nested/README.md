# 03nested/
```console
go run 03nested/main.go --help || :
Usage of /tmp/go-build1880696983/b001/exe/main:
      --another-db.debug        ENV: ANOTHER_DB_DEBUG	-
      --another-db.uri string   ENV: ANOTHER_DB_URI	-
      --db.debug                ENV: DB_DEBUG	-
      --db.uri string           ENV: DB_URI	-
pflag: help requested
exit status 2


$ DB_URI=sqlite:///data.db go run 03nested/main.go --db.debug --another-db.uri sqlite:///:memory:
parsed
{
  "DB": {
    "URI": "sqlite:///data.db",
    "Debug": true
  },
  "AnotherDB": {
    "URI": "sqlite:///:memory:",
    "Debug": false
  }
}
```
