/*
Copyright © 2020 elonzh <elonzh@qq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/elonzh/trumpet/transformers"
	"github.com/sirupsen/logrus"
)

var BuiltinTransformers = []*transformers.Transformer{
	{
		FileName: "builtin/dingtalk-to-feishu.star",
		Src: `
def transform(request):
    msg_type = request["body"]["msgtype"]
    body = {}
    if msg_type == "text":
        body = {
            "msg_type": "text",
            "content": {"text": request["body"]["text"]["content"]},
        }
    elif msg_type == "markdown":
        body = {"msg_type": "interactive", "card": {"elements": []}}
        title = request["body"]["markdown"].get("title")
        if title:
            body["card"]["header"] = {"title": {"content": title, "tag": "plain_text"}}
        text = request["body"]["markdown"].get("text", "")
        body["card"]["elements"].append(
            {"tag": "div", "text": {"content": text, "tag": "lark_md"}}
        )
    request["body"] = body
    return request
`,
	},
	{
		FileName: "builtin/feishu-to-dingtalk.star",
		Src: `
def transform(request):
    msg_type = request["body"]["msg_type"]
    body = {}
    if msg_type == "text":
        body = {
            "msgtype": "text",
            "text": {"content": request["body"]["content"]["text"]},
        }
    request["body"] = body
    return request
`,
	},
}

type Config struct {
	LogLevel        string
	TransformersDir string
	Transformers    []*transformers.Transformer

	allTransformers []*transformers.Transformer
	m               map[string]*transformers.Transformer
}

func (c *Config) LoadAllTransformers() {
	allTransformers := make([]*transformers.Transformer, 0, len(BuiltinTransformers))
	allTransformers = append(allTransformers, BuiltinTransformers...)
	allTransformers = append(allTransformers, cfg.Transformers...)
	if cfg.TransformersDir != "" {
		trans, err := transformers.Load(cfg.TransformersDir)
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		allTransformers = append(allTransformers, trans...)
	}
	c.allTransformers = allTransformers

	m := make(map[string]*transformers.Transformer, len(cfg.Transformers))
	for _, t := range allTransformers {
		if err := t.InitThread(); err != nil {
			logrus.WithError(err).Fatalln()
		}
		if _, exists := m[t.Name]; exists {
			logrus.WithField("Name", t.Name).Warnln("Transformer already exists")
		}
		m[t.Name] = t
	}
	c.m = m
}

func (c *Config) GetTransformer(name string) (*transformers.Transformer, bool) {
	t, ok := c.m[name]
	return t, ok
}
