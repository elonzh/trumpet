package transformers

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkjson"
)

const (
	transformFunctionName = "transform"
)

type Transformer struct {
	Name      string
	src       interface{}
	thread    *starlark.Thread
	transFunc starlark.Value
}

func (t *Transformer) String() string {
	return fmt.Sprintf("Transformer{Name: %s}", t.Name)
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

var namePattern = regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$")

func validateName(s string) error {
	if namePattern.MatchString(s) {
		return nil
	}
	return fmt.Errorf("%s is not a valid url slug as transformer name", s)
}

func NewTransformer(name string, src interface{}) (*Transformer, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	thread := &starlark.Thread{
		Name: name,
	}
	predeclared := starlark.StringDict{
		"json": starlarkjson.Module,
	}
	globals, err := starlark.ExecFile(thread, name, src, predeclared)
	if err != nil {
		return nil, err
	}
	transFunc, ok := globals[transformFunctionName]
	if !ok {
		return nil, fmt.Errorf("transformer not found")
	}
	t := &Transformer{
		Name:      name,
		src:       src,
		thread:    thread,
		transFunc: transFunc,
	}
	return t, nil
}
