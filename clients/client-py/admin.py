import asyncio
import base64

import env
from tcp import TcpClient


async def admin():
    try:
        client = TcpClient()
        await client.init(env.HOST, env.SUPER_ADMIN_PORT)
        client2 = TcpClient()
        await client2.init(env.HOST, 10001)
        client2.write('id-assign')
        r = await client2.read()
        print(r)
        client.write("admin-create " +
                     base64.b64encode(bytes("test", "utf-8")).decode("utf-8") + " " +
                     base64.b64encode(bytes(env.ID, "utf-8")).decode("utf-8"))
        client2.write("set-parameters 0 " + base64.b64encode(bytes("{\"mapSize\":3,\"soldierSpeed\":1,"
                                                                  "\"soldierCreationSpeed\":1,"
                                                                  "\"terrainChangeSpeed\":1,"
                                                                  "\"gameLength\":300}", "utf-8")).decode("utf-8"))
        await client2.read()
        input("Admin create done")
        client.write("start 0")
    except Exception as e:
        print(e)

if __name__ == '__main__':
    asyncio.run(admin())
