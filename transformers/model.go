package transformers

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkjson"
)

const (
	FileSuffix            = ".star"
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

func NewTransformer(filename string, src interface{}) (*Transformer, error) {
	if !strings.HasSuffix(filename, FileSuffix) {
		return nil, fmt.Errorf("filename %s has no suffix %s", filename, FileSuffix)
	}
	name := strings.TrimSuffix(filename, FileSuffix)
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
		Name:      name,
		src:       src,
		thread:    thread,
		transFunc: transFunc,
	}
	return t, nil
}

func Load(dir string) (map[string]*Transformer, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	rv := make(map[string]*Transformer)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), FileSuffix) {
			filename := path.Join(dir, f.Name())
			t, err := NewTransformer(filename, nil)
			if err != nil {
				return nil, fmt.Errorf("error when load %s: %w", filename, err)
			}
			rv[t.Name] = t
		}
	}
	return rv, nil
}
