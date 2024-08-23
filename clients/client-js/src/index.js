const env = require('./env');
const Net = require('node:net');

const client = new Net.Socket();
const Agent = require('./update');

const agent = new Agent();

function action(data) {
    let a = '';
    let byteArray = new Uint8Array(4);
    byteArray[0] = env.MAGIC & 0xFF;
    byteArray[1] = (env.MAGIC >> 8) & 0xFF;
    byteArray[2] = (env.MAGIC >> 16) & 0xFF;
    byteArray[3] = (env.MAGIC >> 24) & 0xFF;

    let length = new Uint8Array(4);
    length[0] = data.length & 0xFF;
    length[1] = (data.length >> 8) & 0xFF;
    length[2] = (data.length >> 16) & 0xFF;
    length[3] = (data.length >> 24) & 0xFF;

    for (let i = 0; i < 4; i++) {
        a += String.fromCharCode(byteArray[i]);
    }

    for (let i = 0; i < 4; i++) {
        a += String.fromCharCode(length[i]);
    }
    a += 'action ' + new Buffer(env.ID).toString('base64') + ' ' + new Buffer(byteArray).toString('base64');
    client.write(a);
}

/**
 * @param msg {string}
 */
function onMessage(msg) {
    /**
     * @type {{type: string, content: string}}
     */
    const obj = JSON.parse(msg);
    if (obj.type === 'action') {
        /**
         * @type {
         *     {
         *         players: {name: string, color: string, numberOfKill: number, possessedTerrainsCount: number}[],
         *         terrains: {terrainId: string, terrainType: number, ownerIndex: number, numberOfSoldier: number, position: [number, number]}[],
         *         pipes: {length: number, first: number, second: number, soldiers: {ownerIndex: number, soldierCount: number, length: number, upward: boolean}[]}[],
         *     }
         * }
         */
        const update = JSON.parse(new Buffer(obj.content, 'base64').toString());
        const ret = agent.update(update);
        const str = JSON.stringify(ret.map(a => a.serialise()));
    }
}

let message = '';

client.connect({port: env.PORT, host: env.HOST}, () => {
    console.log('Connected to the server...');
});

client.on('end', () => {
    console.log('Client is ending');
});

client.on('data', (stream) => {
    message += stream.toString();

    if (message.indexOf('}') !== -1) {
        onMessage(message.substring(8));
        message = '';
    }
});