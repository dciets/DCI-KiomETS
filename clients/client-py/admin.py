import asyncio
import base64

import env
from tcp import TcpClient


async def admin():
    try:
        client = TcpClient()
        await client.init(env.HOST, env.SUPER_ADMIN_PORT)
        client.write("admin-create " +
                     base64.b64encode(bytes("test", "utf-8")).decode("utf-8") + " " +
                     base64.b64encode(bytes(env.ID, "utf-8")).decode("utf-8"))
        input("Admin create done")
        client.write("start")
    except Exception as e:
        print(e)

if __name__ == '__main__':
    asyncio.run(admin())
