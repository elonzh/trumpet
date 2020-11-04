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
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  ``,
}

type Config struct {
	LogLevel     logrus.Level
	Transformers map[string]*transformers.Transformer
}

var (
	cfg = &Config{
		LogLevel:     logrus.InfoLevel,
		Transformers: map[string]*transformers.Transformer{},
	}
)

func registerTransformers(m map[string]string) {
	for name, src := range m {
		t, err := transformers.NewTransformer(name, src)
		if err != nil {
			logrus.WithError(err).WithField("Name", name).Fatalln("error when init Transformer")
		}
		if _, exists := cfg.Transformers[t.Name]; exists {
			logrus.WithField("Name", name).Warnln("Transformer already exists")
		}
		cfg.Transformers[t.Name] = t
	}
}

func init() {
	builtinTransformers := map[string]string{
		"feishu-to-dingtalk": `
def transform(request):
	msg_type = request["body"]["msg_type"]
	body = {}
	if msg_type == "text":
		body = {"msgtype": "text", "text": {"content": request["body"]["content"]["text"]}}
	request["body"] = body
	return request
`,
		"dingtalk-to-feishu": `
def transform(request):
	msg_type = request["body"]["msgtype"]
	body = {}
	if msg_type == "text":
		body = {"msg_type": "text", "content": {"text": request["body"]["text"]["content"]}}
	elif msg_type == "markdown":
		title = request["body"]["markdown"].get("title")
		text = request["body"]["markdown"].get("text", "")
		if title:
			text = title + "\n" + text
		body = {"msg_type": "text", "content": {"text": text}}
	request["body"] = body
	return request
`,
	}
	registerTransformers(builtinTransformers)
	rootCmd.AddCommand(configCmd)
}
