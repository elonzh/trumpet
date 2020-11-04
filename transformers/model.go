package transformers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkjson"
	"go.starlark.net/starlarkstruct"
)

const (
	transformFunctionName = "transform"
	headerKey             = starlark.String("header")
	bodyKey               = starlark.String("body")
)

var starJsonDecode, starJsonEncode *starlark.Builtin

func init() {
	starDecode, err := starlarkjson.Module.Attr("decode")
	if err != nil {
		panic(err)
	}
	starJsonDecode = starDecode.(*starlark.Builtin)
	starEncode, err := starlarkjson.Module.Attr("encode")
	if err != nil {
		panic(err)
	}
	starJsonEncode = starEncode.(*starlark.Builtin)
}

type Transformer struct {
	Name      string
	src       interface{}
	thread    *starlark.Thread
	transFunc starlark.Value
}

func (t *Transformer) String() string {
	return fmt.Sprintf("Transformer{Name: %s}", t.Name)
}

func (t *Transformer) requestToStarDict(req *http.Request) (*starlark.Dict, error) {
	starHeader := starlark.NewDict(len(req.Header))
	for k, v := range req.Header {
		if len(v) == 1 {
			err := starHeader.SetKey(starlark.String(k), starlark.String(v[0]))
			if err != nil {
				return nil, err
			}
		} else if len(v) > 1 {
			logrus.WithField("Header", req.Header).Warningf("header %s has multiple values", k)
		}
	}
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	starBody, err := starJsonDecode.CallInternal(t.thread, starlark.Tuple([]starlark.Value{starlark.String(rawBody)}), nil)
	if err != nil {
		return nil, err
	}
	starReq := starlark.NewDict(2)
	err = starReq.SetKey(headerKey, starHeader)
	if err != nil {
		return nil, err
	}
	err = starReq.SetKey(bodyKey, starBody)
	if err != nil {
		return nil, err
	}
	return starReq, nil
}

func (t *Transformer) updateRequestFromStarDict(req *http.Request, data *starlark.Dict) error {
	v, found, err := data.Get(headerKey)
	if !found || err != nil {
		return fmt.Errorf("error when get `%s` from StarDict, found %t, error %w", headerKey, found, err)
	}
	newHeader := http.Header{}
	for _, item := range v.(*starlark.Dict).Items() {
		newHeader.Set(string(item.Index(0).(starlark.String)), string(item.Index(1).(starlark.String)))
	}
	req.Header = newHeader
	v, found, err = data.Get(bodyKey)
	if !found || err != nil {
		return fmt.Errorf("error when get `%s` from StarDict, found %t, error %w", bodyKey, found, err)
	}
	encodedBody, err := starJsonEncode.CallInternal(t.thread, starlark.Tuple([]starlark.Value{v}), nil)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(strings.NewReader(string(encodedBody.(starlark.String))))
	return nil
}

func (t *Transformer) Exec(req *http.Request) (*http.Request, error) {
	newReq := req.Clone(req.Context())
	starReq, err := t.requestToStarDict(req)
	if err != nil {
		return nil, err
	}
	result, err := starlark.Call(t.thread, t.transFunc, starlark.Tuple{starReq}, nil)
	if err != nil {
		return nil, err
	}
	starNewReq, ok := result.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("can not convert result: %s", result)
	}
	if err := t.updateRequestFromStarDict(newReq, starNewReq); err != nil {
		return nil, err
	}
	return newReq, nil
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
		starlarkjson.Module.Name: starlarkjson.Module,
		"struct":                 starlark.NewBuiltin("struct", starlarkstruct.Make),
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
