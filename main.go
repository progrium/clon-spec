package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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
	lastIdx := len(parts) - 1
	last := parts[lastIdx]
	for i, next := range parts[1:lastIdx] {
		full = full + "/" + next
		if _, ok := seen[full]; !ok {
			// new hierarchy level
			seen[full] = 0
			if _, err := strconv.Atoi(last); err == nil && i == lastIdx-2 {
				ret = append(ret, addRaw(full, []byte("[]")))
			} else {
				ret = append(ret, addRaw(full, []byte("{}")))
			}
		}
	}
	return ret
}

func translatePath(orig string) (string, bool) {
	isRaw := orig[len(orig)-1:] == ":"
	if isRaw {
		orig = strings.TrimRight(orig, ":")
	}
	trans := ""
	for _, part := range strings.Split(orig, "[") {
		part = strings.TrimRight(part, "]")
		//log.Println("  next part:", part)
		if part == "" {
			// array
			if _, ok := seen[trans]; !ok {
				part = "0"
			} else {
				seen[trans] += 1
				part = strconv.Itoa(seen[trans])
			}
		}
		trans = trans + "/" + part
	}
	return trans, isRaw
}
func toString(op jsonpatch.Operation) string {
	path, _ := op.Path()
	val, _ := op.ValueInterface()
	return fmt.Sprintf("op:%s path:%s val:%s", op.Kind(), path, val)
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
		log.Println("jsonpath:", trPath)
		if raw {
			nextOp = addRaw(trPath, []byte(val))
		} else {
			nextOp = addString(trPath, val)
		}
		for _, op := range getHierarchy(nextOp) {
			// missing from hierarchy
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
