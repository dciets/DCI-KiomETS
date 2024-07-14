import asyncio
import logging
from asyncio import StreamReader, StreamWriter

import env


class TcpClient:
    __reader: StreamReader
    __writer: StreamWriter

    __initialized: bool

    def __init__(self):
        self.__initialized = False
        self.__reader = None
        self.__writer = None

    async def init(self, host: str, port: int):
        self.__initialized = True
        self.__reader, self.__writer = await asyncio.open_connection(host, port)

    async def read(self) -> tuple[str, bool]:
        if not self.__initialized:
            return "", False
        try:
            read: bytes = await self.__reader.read(8)
            if len(read) != 8:
                return "", False
            magic: int = int.from_bytes(read[0:4], 'little')
            length: int = int.from_bytes(read[4:8], 'little')

            if magic != env.MAGIC:
                return "", False

            read: bytes = await self.__reader.read(length)
            return str(read, 'utf-8'), True
        except ConnectionError as err:
            self.__initialized = False
            self.__writer.close()
            self.__writer = None
            self.__reader = None
            print(err)
            return "", False

    def write(self, content: str) -> bool:
        if not self.__initialized:
            return False
        try:
            magic: bytes = int.to_bytes(env.MAGIC, 4, 'little')
            length: bytes = int.to_bytes(len(content), 4, 'little')

            self.__writer.write(magic + length + bytes(content, 'utf-8'))
            return True
        except ConnectionError as err:
            self.__initialized = False
            self.__writer.close()
            self.__writer = None
            self.__reader = None
            print(err)
            return False
