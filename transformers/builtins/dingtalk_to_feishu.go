package builtins

import (
	"log"

	"github.com/elonzh/trumpet/transformers"
)

var dingtalkToFeishu *transformers.Transformer

func init() {
	name := "dingtalk_to_feishu.star"
	src := `
def transform(raw):
	origin_body = json.decode(raw)
	msg_type = origin_body['msgtype']
	body = {}
	if msg_type == "text":
		body = {"msg_type": "text", "content": {"text": origin_body["text"]["content"]}}
	elif msg_type == "markdown":
		title = origin_body["markdown"].get("title")
		text = origin_body["markdown"].get("text", "")
		if title:
			text = title + "\n" + text
		body = {"msg_type": "text", "content": {"text": text}}
	return json.encode(body)
`
	dingtalkToFeishu, err := transformers.NewTransformer(name, src)
	if err != nil {
		log.Fatalf("error when init Transformer %s, %s", name, err)
	}
	register(dingtalkToFeishu)
}
