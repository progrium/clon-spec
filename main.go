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

var seen map[string]int = map[string]int{}

func getHierarchy(op jsonpatch.Operation) []jsonpatch.Operation {
	ret := []jsonpatch.Operation{}
	p, _ := op.Path()
	parts := strings.Split(p, "/")
	full := ""
	for _, next := range parts[1 : len(parts)-1] {
		full = full + "/" + next
		if _, ok := seen[full]; !ok {
			seen[full] = 0
			ret = append(ret, addRaw(full, []byte("{}")))
		}
	}
	return ret
}

func translatePath(orig string) (string, bool) {
	isRaw := orig[len(orig)-1:] == ":"
	trans := orig
	if isRaw {
		trans = strings.TrimRight(orig, ":")
	}
	trans = strings.ReplaceAll(trans, ".", "/")
	return "/" + trans, isRaw
}

func parseArgs() jsonpatch.Patch {
	patch := jsonpatch.Patch{}
	for i, next := range os.Args[1:] {
		log.Printf("arg[%d]: %q\n", i, next)
		parts := strings.Split(next, "=")
		path, val := parts[0], parts[1]

		var nextOp jsonpatch.Operation
		trPath, raw := translatePath(path)
		seen[trPath] = 1
		log.Println("trPat:", trPath)
		if raw {
			nextOp = addRaw(trPath, []byte(val))
		} else {
			nextOp = addString(trPath, val)
		}
		for _, op := range getHierarchy(nextOp) {
			patch = append(patch, op)
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
