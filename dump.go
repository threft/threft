package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/threft/threft/tidm"
	"io"
	"os"
)

func dumpTIDM(t *tidm.TIDM) error {
	dumpFile, err := os.Create("tidm_dump")
	if err != nil {
		return fmt.Errorf("Error creating dumpfile: %s", err)
	}

	// write spew dump
	io.WriteString(dumpFile, "spew.Dump:\n==========\n")
	cs := spew.NewDefaultConfig()
	cs.Indent = "    "
	cs.Fdump(dumpFile, t)
	io.WriteString(dumpFile, "\n\n\n\n")

	// write json
	io.WriteString(dumpFile, "tidm-json:\n==========\n")
	enc := json.NewEncoder(dumpFile)
	err = enc.Encode(t)
	if err != nil {
		return fmt.Errorf("Error encoding tidm-json to dumpfile: %s", err)
	}

	// all done
	return nil
}
