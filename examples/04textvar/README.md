# 04textvar/
```console
go run 04textvar/main.go --help || :
Usage of textvar:
      --ip *net.IP   ENV: X_IP	-
      --ip2 net.IP   ENV: X_IP2	- (default 0.0.0.1)
pflag: help requested
exit status 2


$ go run 04textvar/main.go --ip 127.0.0.1
parsed: &main.Options{IP:(*net.IP)(0xc0000a6078), IP2:net.IP{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0x0, 0x1}}
-------------------------------------

$ json:
{
  "IP": "127.0.0.1",
  "IP2": "0.0.0.1"
}
```
