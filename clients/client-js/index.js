import { env } from './env.js';
import { Socket } from 'node:net';
import { ArgumentParser } from 'argparse';
import Agent from './update.js';

const client = new Socket();


const agent = new Agent();

function action(data) {
    let byteArray = new Uint8Array(4);
    byteArray[0] = env.magic & 0xFF;
    byteArray[1] = (env.magic >> 8) & 0xFF;
    byteArray[2] = (env.magic >> 16) & 0xFF;
    byteArray[3] = (env.magic >> 24) & 0xFF;

    let msg = ""

    for (let i = 0; i < 4; i++) {
        msg += String.fromCharCode(byteArray[i]);
    }
    let action = 'action ' + new Buffer(env.id).toString('base64') + ' ' + new Buffer(data).toString('base64');

    let length = new Uint8Array(4);
    let actionLength = action.length
    length[0] = actionLength & 0xFF;
    length[1] = (actionLength >> 8) & 0xFF;
    length[2] = (actionLength >> 16) & 0xFF;
    length[3] = (actionLength >> 24) & 0xFF;

    for (let i = 0; i < 4; i++) {
        msg += String.fromCharCode(length[i]);
    }
    client.write(msg + action);
}

/**
 * @param msg {UInt8Array}
 */
function onMessage(msg) {
    /**
     * @type {{type: string, content: string}}
     */
    const header = msg.slice(0, 8);
    msg = msg.slice(8);

    const magic = header[0] | (header[1] << 8) | (header[2] << 16) | (header[3] << 24);
    if (magic !== env.magic) {
        console.error('Invalid magic number:', magic);
        return;
    }
    const length = header[4] | (header[5] << 8) | (header[6] << 16) | (header[7] << 24);
    if (length !== msg.length) {
        console.error('Invalid message length:', length, 'expected:', msg.length);
        return;
    }
    const decoder = new TextDecoder('utf-8');
    const str = decoder.decode(msg);
    const obj = JSON.parse(str);
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
        action(str);
    }
}



let parser = ArgumentParser({ description: "TCP Client for the game" })
parser.add_argument("bot_id", { type: "str", help: "ID of the bot" })
parser.add_argument("bot_name", { type: "str", help: "Name of the bot" })
let parsed_args = parser.parse_args();

env.id = parsed_args.bot_id
env.playerName = parsed_args.bot_name


let message = Uint8Array.from([]);

client.connect({ port: env.port, host: env.host }, () => {
    console.log('Connected to the server...');
});

client.on('end', () => {
    console.log('Client is ending');
});

client.on('data', (stream) => {
    message = new Uint8Array([...message, ...stream]);

    if (message.indexOf(0x7d) !== -1) {
        onMessage(message);
        message = Uint8Array.from([]);
    }
});