import asyncio
import base64
import json
from typing import Tuple

import env
import server
from actions import MoveAction, Action
from connection import Message
from tcp import TcpClient
from update import Agent


def connect(client: TcpClient):
    client.write("connect " + base64.b64encode(bytes(env.ID, 'utf-8')).decode('utf-8'))


def action(client: TcpClient, actions_json: str):
    client.write("action " + base64.b64encode(bytes(env.ID, 'utf-8')).decode('utf-8') + " " + base64.b64encode(bytes(actions_json, 'utf-8')).decode('utf-8'))


async def tcp_client():
    try:
        client = TcpClient()
        await client.init(env.HOST, env.PORT)
        bot = Agent()
        while True:
            json_data: Tuple[str, bool] = await client.read()
            json_string: str = json_data[0]
            msg = Message(json.loads(json_string))
            if msg.type() == "action":
                print(msg.content())
                game = server.convert_data(msg.content())
                actions = bot.update(game)
                actions_json = json.dumps([ac.serialize() for ac in actions])
                action(client, actions_json)

    except Exception as e:
        print(e)


if __name__ == "__main__":
    asyncio.run(tcp_client())
