package transformers

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

const suffix = ".star"

var transformers = map[string]*Transformer{}

func Register(t *Transformer) {
	o, exists := transformers[t.name]
	if exists {
		log.Printf("%s already exists, original data will be replaced", o)
	}
	transformers[t.name] = t
}

func Get(name string) (*Transformer, bool) {
	t, ok := transformers[name]
	return t, ok
}

func Load(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), suffix) {
			filename := path.Join(dir, f.Name())
			t, err := NewTransformer(filename, nil)
			if err != nil {
				return fmt.Errorf("error when load %s: %w", filename, err)
			}
			Register(t)
		}
	}
	return nil
}

func init() {
	for name, code := range map[string]string{
		"dingtalk_to_feishu.star": `
def transform(raw):
	origin_body = json.decode(raw)
	msg_type = origin_body['msgtype']
	body = {}
	if msg_type == "text":
		body = {"msg_type": "text", "content": {"text": origin_body["text"]["content"]}}
	return json.encode(body)
`,
		"feishu_to_dingtalk.star": `
def transform(raw):
	origin_body = json.decode(raw)
	msg_type = origin_body['msg_type']
	body = {}
	if msg_type == "text":
		body = {"msgtype": "text", "text": {"content": origin_body["content"]["text"]}}
	return json.encode(body)
`,
	} {
		t, err := NewTransformer(name, code)
		if err != nil {
			log.Fatalln("error when init transformers,", err)
		}
		transformers[t.name] = t
	}
}
