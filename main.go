package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func raw(s string) *json.RawMessage {
	r := json.RawMessage(fmt.Sprintf("%q", s))
	return &r
}
func addString(path, val string) jsonpatch.Operation {
	return jsonpatch.Operation{
		"op":    raw("add"),
		"path":  raw(path),
		"value": raw(val),
	}
}
func addRaw(path string, val []byte) jsonpatch.Operation {
	v := json.RawMessage(val)
	return jsonpatch.Operation{
		"op":    raw("add"),
		"path":  raw(path),
		"value": &v,
	}
}
func parseArgs() jsonpatch.Patch {
	patch := jsonpatch.Patch{}
	for i, next := range os.Args[1:] {
		log.Printf("arg[%d]: %q\n", i, next)
		parts := strings.Split(next, "=")
		path, val := parts[0], parts[1]

		var nextOp jsonpatch.Operation
		if path[len(path)-1:] == ":" {
			trimmedPath := strings.TrimRight(path, ":")
			nextOp = addRaw("/"+trimmedPath, []byte(val))
		} else {
			nextOp = addString("/"+path, val)
		}
		patch = append(patch, nextOp)
	}
	return patch
}

func main() {
	log.SetFlags(log.Lmsgprefix)
	patch := parseArgs()

	modified, err := patch.Apply([]byte(`{}`))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(modified))
}
