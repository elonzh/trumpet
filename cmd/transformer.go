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
	"fmt"

	"github.com/elonzh/trumpet/transformers"
	"github.com/spf13/cobra"
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

func newTransformerCmd(cfg *Config) *cobra.Command {
	transformerCmd := &cobra.Command{
		Use: "transformer",
	}
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			for _, t := range cfg.GetAllTransformers() {
				fmt.Println(t.Sprint())
			}
		},
	}
	transformerCmd.AddCommand(listCmd)
	return transformerCmd
}
