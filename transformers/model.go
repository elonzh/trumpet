package transformers

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkjson"
)

const transformFunctionName = "transform"

type Transformer struct {
	name      string
	src       interface{}
	thread    *starlark.Thread
	transFunc starlark.Value
}

func (t *Transformer) String() string {
	return fmt.Sprintf("Transformer{name: %s}", t.name)
}

func (t *Transformer) Exec(raw string) (string, error) {
	args := starlark.Tuple{starlark.String(raw)}
	result, err := starlark.Call(t.thread, t.transFunc, args, nil)
	if err != nil {
		return "", err
	}
	rv, ok := starlark.AsString(result)
	if !ok {
		return "", fmt.Errorf("can not convert result as string: %s", rv)
	}
	return rv, nil
}

func NewTransformer(filename string, src interface{}) (*Transformer, error) {
	if !strings.HasSuffix(filename, suffix) {
		return nil, fmt.Errorf("filename %s has no suffix %s", filename, suffix)
	}
	name := strings.TrimSuffix(filename, suffix)
	thread := &starlark.Thread{
		Name: name,
	}
	predeclared := starlark.StringDict{
		"json": starlarkjson.Module,
	}
	globals, err := starlark.ExecFile(thread, filename, src, predeclared)
	if err != nil {
		return nil, err
	}
	transFunc, ok := globals[transformFunctionName]
	if !ok {
		return nil, fmt.Errorf("transformer not found")
	}
	t := &Transformer{
		name:      name,
		src:       src,
		thread:    thread,
		transFunc: transFunc,
	}
	return t, nil
}
