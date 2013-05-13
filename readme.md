## Threft application

The threft application is the user front-end for threft code generation. It parses .thrift files and either generates code or invokes an other generator executable which generates the code.

TODO: Test how submodule behaves in this.. i.e. submodule tidm here, and in gen-go, does that work? Probably wont work as expected and result in lots of trouble.

Imports gen-go, tidm, encoding/json(no it wont), flags(?).

**parser/generator communication**
For now: skip the marshalling/unmarshalling part. Directly invoke the generator with direct access to the TIDM structure.

Future: Clearly describe the difference between "direct" generators (gen-go, which can be included in the threft binary) and "indirect" generators (coupled with marshalling/unmarshalling the TIDM structure)

For marshalling/unmarshalling: consider rjson (readable json):
http://rogpeppe.wordpress.com/2012/09/24/goson-readable-json/
http://go.pkgdoc.org/launchpad.net/rjson
for rjson: test speed and stability, as well as functionality compared to encoding/json

For marshalling/unmarshalling: consider using threft as communication for the parser/generator. Parser becomes server, generator client. using stdin/out on the client process as transport. This way the generator receives only the scope it is interested in. I prefer this idea.
