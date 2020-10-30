package builtins

import (
	"log"

	"github.com/elonzh/trumpet/transformers"
)

var feishuToDingtalk *transformers.Transformer

func init() {
	name := "feishu_to_dingtalk.star"
	src := `
def transform(raw):
	origin_body = json.decode(raw)
	msg_type = origin_body['msg_type']
	body = {}
	if msg_type == "text":
		body = {"msgtype": "text", "text": {"content": origin_body["content"]["text"]}}
	return json.encode(body)
`
	feishuToDingtalk, err := transformers.NewTransformer(name, src)
	if err != nil {
		log.Fatalf("error when init Transformer %s, %s", name, err)
	}
	register(feishuToDingtalk)
}
