## Threft

The threft application is the user front-end for threft code generation. It parses .thrift files and invokes a generator which generates code.

Marshalling is done with tidm-json

For marshalling/unmarshalling: consider rjson (readable json):
http://rogpeppe.wordpress.com/2012/09/24/goson-readable-json/
http://go.pkgdoc.org/launchpad.net/rjson
for rjson: test speed and stability, as well as functionality compared to encoding/json