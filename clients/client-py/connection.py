import base64
import json


class Message:
    __type: str
    __content: str

    def __init__(self, json_str):
        self.__type = json_str['type']
        self.__content = base64.b64decode(json_str['content']).decode('utf-8')

    def type(self) -> str:
        return self.__type

    def content(self) -> str:
        return self.__content
