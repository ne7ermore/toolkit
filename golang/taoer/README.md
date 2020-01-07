## taoer: A Golang for Pressure Test.

simplely and human-friendly

### Installation/update
```
go get -u github.com/ne7ermore/taoer
```

### Use
```
step-1: go build github.com/ne7ermore/taoer
step-2: ./taoer -url=http://example -qps=100 -duration=60
```

### More info:
```
The following arguments are mandatory:
  -url               server url
  -form              http form eg: 'a=a&b=b' required if not GET

The following arguments are optional:
  -method            http method ['GET']
  -duration          seconds for requst   ['60']
  -qps               requst per second ['100'] note: 10 multiples
  -disKA             DisableKeepAlives, if true, prevents re-use of TCP connections between different HTTP requests ['true']
  -disComp           DisableCompression, if true, prevents the Transport from requesting compression with an Accept-Encoding: gzip ['true']
  -timeout           time for handshake ['0']
  -cpus              number of CPUs ['maximum']
```
