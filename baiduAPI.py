import base64

import requests
import os

class BaiduOCR(object):
    """docstring for ClassName.

    Attributes:

    """

    def __init__(self):
        self.get_token()
    def get_APIKEY(self):
        if os.path.exists("api.txt"):
            with open("api.txt", "r") as f:
                line = f.readline()
                self.key = line.split(':')
        else:
            print("push.txt NOT FOUND. Initialising for token")
            with open("api.txt", "w") as f:
                API_key = input("API_key:")
                secret_key = input("secret_key:")
                self.key = (API_key, secret_key)
                f.writelines(':'.join(self.key))
        return self.key
    def get_token(self):
        self.get_APIKEY()
        _id = self.key[0]
        _secret = self.key[1]

        url = "https://aip.baidubce.com/oauth/2.0/token"
        url += f'?grant_type=client_credentials&client_id={_id}&client_secret={_secret}'

        response = requests.get(url)
        res = response.json()

        self.access_token = res['access_token']
        return self.access_token

    def recognize(self, img):
        '''
        通用文字识别（高精度版）
        '''
        self.get_token()
        headers = {'content-type': 'application/x-www-form-urlencoded'}

        request_url = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
        request_url = request_url + "?access_token=" + self.access_token

        params = {"image": img}
        response = requests.post(request_url, data=params, headers=headers)

        if not response.status_code == 200:
            return 'error'

        res = response.json()
        if not res.get('words_result_num'):
            print(res)
            return 'error'

        return res['words_result'][0]['words'].replace(" ", "")

    def recognize_general_basic(self, img):
        '''
        通用文字识别（标准版）
        '''
        self.get_token()
        headers = {'content-type': 'application/x-www-form-urlencoded'}

        request_url = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
        request_url = request_url + "?access_token=" + self.access_token

        params = {"image": img}
        response = requests.post(request_url, data=params, headers=headers)

        if not response.status_code == 200:
            return 'error'

        res = response.json()
        if not res.get('words_result_num'):
            print(res)
            return 'error'

        return res['words_result'][0]['words'].replace(" ", "")


def test(folder, fudan, ocr):
    code_img = fudan.get_code_img()
    code_text = ocr.recognize(code_img)
    gen_code_text = ocr.recognize_general_basic(code_img)

    print('ocr', code_text, gen_code_text)

    with open(f"{folder}/accuracy_{code_text}.png", "wb") as fh:
        fh.write(base64.decodebytes(code_img))

    with open(f"{folder}/general_{gen_code_text}.png", "wb") as fh:
        fh.write(base64.decodebytes(code_img))

# if __name__ == "__main__":
#
#     import os
#     import sys
#     import base64
#     from dailyFudan import get_config, Zlapp
#
#     config = get_config('./config.yml')
#
#     service = config['service']
#
#     ocr = BaiduOCR()
#
#     if not ocr.get_token(service['ocr']):
#         sys.exit()
#
#     daily_fudan = Zlapp("user['uid']", "user['psw']", url_login="zlapp_login")
#
#     folder = './ocrTest'
#
#     if not os.path.exists(folder):
#         os.makedirs(folder)
#
#     for i in range(10):
#         test(folder, daily_fudan, ocr)