import os
import unittest

import requests

feishu_webhook = os.getenv("TRUMPET_FEISHU_WEBHOOK")
dingtalk_webhook = os.getenv("TRUMPET_DINGTALK_WEBHOOK")

# https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
# https://www.feishu.cn/hc/zh-cn/articles/360024984973-%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%99%A8%E4%BA%BA


class TestStringMethods(unittest.TestCase):
    def test_dingtalk_to_feishu(self):
        url = f"http://127.0.0.1:8080/transformers/dingtalk_to_feishu?trumpet_to={feishu_webhook}"
        cases = [
            {"msgtype": "text", "text": {"content": "快乐小神仙"}},
            {
                "markdown": {
                    "title": "哈哈哈 触发了 job test, 构建号：620",
                    "text": "###### 项目 [Unob](https://coding.net/p/unob)\n[哈哈哈](https://coding.net/u/ljaSkNTntD) 触发了 job \n> [test](https://coding.net/p/unob/ci/job/260491) 构建号：[620](https://coding.net/p/proj/ci/job/260491/build/620/pipeline)",
                },
                "msgtype": "markdown",
            },
            {
                "markdown": {
                    "text": "###### 项目 [Unob](https://coding.net/p/unob)\n[哈哈哈](https://coding.net/u/ljaSkNTntD) 触发了 job \n> [test](https://coding.net/p/unob/ci/job/260491) 构建号：[620](https://coding.net/p/proj/ci/job/260491/build/620/pipeline)",
                },
                "msgtype": "markdown",
            },
        ]
        for case in cases:
            resp = requests.post(url, json=case)
            print(resp.status_code)
            print(resp.headers)
            print(resp.text)
            assert resp.ok

    def test_feishu_to_dingtalk(self):
        url = f"http://127.0.0.1:8080/transformers/feishu_to_dingtalk?trumpet_to={dingtalk_webhook}"
        cases = [{"msg_type": "text", "content": {"text": "快乐小神仙"}}]
        for case in cases:
            resp = requests.post(url, json=case)
            print(resp.status_code)
            print(resp.headers)
            print(resp.text)
            assert resp.ok


if __name__ == "__main__":
    unittest.main()
