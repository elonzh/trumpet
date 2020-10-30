package builtins

import (
	"log"

	"github.com/elonzh/trumpet/transformers"
)

var ts = map[string]*transformers.Transformer{}

func register(t *transformers.Transformer) {
	o, exists := ts[t.Name]
	if exists {
		log.Fatalf("%s already exists", o)
	}
	ts[t.Name] = t
}

func Get(name string) (*transformers.Transformer, bool) {
	t, ok := ts[name]
	return t, ok
}
