import os
import requests

feishu_webhook = os.getenv("TRUMPET_FEISHU_WEBHOOK")
dingtalk_webhook = os.getenv("TRUMPET_DINGTALK_WEBHOOK")

# https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
# https://www.feishu.cn/hc/zh-cn/articles/360024984973-%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%99%A8%E4%BA%BA
resp = requests.post(
    f"http://127.0.0.1:8080/transformers/dingtalk_to_feishu?trumpet_to={feishu_webhook}",
    json={"msgtype": "text", "text": {"content": "快乐小神仙"}}
)
print(resp.status_code)
print(resp.headers)
print(resp.text)

resp = requests.post(
    f"http://127.0.0.1:8080/transformers/feishu_to_dingtalk?trumpet_to={dingtalk_webhook}",
    json={"msg_type": "text", "content": {"text": "快乐小神仙"}}
)
print(resp.status_code)
print(resp.headers)
print(resp.text)
