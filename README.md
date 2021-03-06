# 🎺trumpet ![GitHub release (latest by date)](https://img.shields.io/github/v/release/elonzh/trumpet?style=flat-square) ![Docker Pulls](https://img.shields.io/docker/pulls/elonzh/trumpet?style=flat-square) [![GolangCI](https://golangci.com/badges/github.com/elonzh/trumpet.svg)](https://golangci.com) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/elonzh/trumpet/build?style=flat-square) [![GitHub license](https://img.shields.io/github/license/elonzh/trumpet?style=flat-square)](https://github.com/elonzh/trumpet/blob/main/LICENSE)

Webhook message transform service

---

[English](./README.md) | [简体中文](./README.zh.md)

## Usage

### Quick start

```shell
docker run -d -p 8080:8080 elonzh/trumpet
feishu_webhook='<your-feishu-webhook-url>'
curl "http://127.0.0.1:8080/transformers/dingtalk-to-feishu?trumpet_to=${feishu_webhook}" \
    -v -X "POST" -H "Content-Type: application/json" \
    -d '{"msgtype": "text", "text": {"content": "message from trumpet!"}}'
```

You can mount the configuration in the default configuration path `/app/config.yaml`, or provide the `-c/--config` parameter to provide the configuration file path.

### Builtin transformers

|                                                                  |   |                                                                                                                                                                                                                                                         |
|------------------------------------------------------------------|---|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [DingTalk](https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq) | ↔ | [Feishu](https://www.feishu.cn/hc/zh-cn/articles/360024984973-%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%99%A8%E4%BA%BA)/[Lark](https://www.larksuite.com/hc/en-US/articles/360048487736-Bot-Use-bots-in-groups#source=section) |

### Customize transformers

> Starlark is a dialect of Python intended for use as a configuration language.

The message transformer are written by [Starlark language](https://github.com/google/starlark-go),
and what you need to do is defining a `transform` function, modifying the incoming request accordingly,
for example:

```python
def transform(request):
    # print(requst["headers"])
    # print(requst["body"])
    msg_type = request["body"]["msg_type"]
    body = {}
    if msg_type == "text":
        body = {
            "msgtype": "text",
            "text": {"content": request["body"]["content"]["text"]},
        }
    request["body"] = body
    return request
```

### Deploy to Kubernetes

The configuration file deployed to the Kubernetes cluster is provided in the [manifests](./manifests) folder, you can make adjustments according to your needs.

## Contribute

If you want to add or update builtin transformers, just make a pull request!
