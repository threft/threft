`tidm` is an acronym for Threft Interface Definition Model.
It contains methods to parse a .threft file to a in-memory interface definition model and provides marshalling to `tidm-json`.

Does not marshall to JSON.

// Parses a .threft IDL file to TisObject
ParseThreftFile(filename string) (*tidm.TIDM, error)


Parsing:

1) A list of files is given, these files are saved in-memory as Documents.
2) Each Document is parsed, which creates one or more namespaces.
2a) For each document, a namespace is created in the 'default' Target.
2b) For each 'namespace header' in a document, a namespace is created in the Target as defined by the 'namespacescope' value in the 'namespace header'
