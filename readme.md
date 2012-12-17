## Threft application

The threft application is the user front-end for threft code generation. It allows .thrift parsing and code generation for officially supported code generators.

TODO: Test how submodule behaves in this.. i.e. submodule tidm here, and in gen-go, does that work? Probably wont work as expected and deliver troubles.

Imports gen-go(no it wont), tidm, encoding/json, flags(?).

Consider rjson (readable json):
http://rogpeppe.wordpress.com/2012/09/24/goson-readable-json/
http://go.pkgdoc.org/launchpad.net/rjson
for rjson: test speed and stability, as well as functionality compared to encoding/json