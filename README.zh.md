# 🎺trumpet ![GitHub release (latest by date)](https://img.shields.io/github/v/release/elonzh/trumpet?style=flat-square) ![Docker Pulls](https://img.shields.io/docker/pulls/elonzh/trumpet?style=flat-square) [![GolangCI](https://golangci.com/badges/github.com/elonzh/trumpet.svg)](https://golangci.com) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/elonzh/trumpet/build?style=flat-square) [![GitHub license](https://img.shields.io/github/license/elonzh/trumpet?style=flat-square)](https://github.com/elonzh/trumpet/blob/main/LICENSE)

Webhook 消息转换服务

---

[English](./README.md) | [简体中文](./README.zh.md)

## 使用

### 快速上手

```shell
docker run -d -p 8080:8080 elonzh/trumpet
feishu_webhook='<your-feishu-webhook-url>'
curl "http://127.0.0.1:8080/transformers/dingtalk-to-feishu?trumpet_to=${feishu_webhook}" \
    -v -X "POST" -H "Content-Type: application/json" \
    -d '{"msgtype": "text", "text": {"content": "message from trumpet!"}}'
```

你可以将自定义配置挂载在默认的配置路径 `/app/config.yaml` ，或者提供 `-c/--config` 参数提供配置文件路径。

### 内置的消息转换器

|                                                                  |   |                                                                                                                                                                                                                                                         |
|------------------------------------------------------------------|---|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [DingTalk](https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq) | ↔ | [Feishu](https://www.feishu.cn/hc/zh-cn/articles/360024984973-%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%99%A8%E4%BA%BA)/[Lark](https://www.larksuite.com/hc/en-US/articles/360048487736-Bot-Use-bots-in-groups#source=section) |

### 自定义消息转换器

> Starlark is a dialect of Python intended for use as a configuration language.

自定义消息转换器只需要用 [Starlark 语言](https://github.com/google/starlark-go) 定义一个 `transform` 函数，对传入的请求做相应的修改即可，例如：

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

### 部署至 Kubernetes 集群

[manifests](./manifests) 文件夹内提供了部署至 Kubernetes 集群的配置文件，你可根据需求做相应的调整。

## 贡献代码

你可以提交合并请求更新或添加内置的转换器供大家使用。
