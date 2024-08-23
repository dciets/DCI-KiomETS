const actions = require('./actions');

module.exports = class Agent {
    /**
     *
     * @param data {
     *     {
     *         players: {name: string, color: string, numberOfKill: number, possessedTerrainsCount: number}[],
     *         terrains: {terrainId: string, terrainType: number, ownerIndex: number, numberOfSoldier: number, position: [number, number]}[],
     *         pipes: {length: number, first: number, second: number, soldiers: {ownerIndex: number, soldierCount: number, length: number, upward: boolean}[]}[],
     *     }
     * }
     */
    update(data) {
        const playerName = 'test';

        const terrains = data.terrains;
        const players = data.players;
        const pipes = data.pipes;

        const player = players.find(p => p.name === playerName);
        if (player) {
            const playerIndex = players.indexOf(player);
            const playerTerrains = terrains.filter(t => t.name === playerIndex);

            const playerTerrainsIndex = playerTerrains.map(t => terrains.indexOf(t));

            const orders = [];
            for (const pipe of pipes) {
                if (playerTerrainsIndex.indexOf(pipe.first) !== -1 && playerTerrainsIndex.indexOf(pipe.second) === -1 && terrains[pipe.first].numberOfSoldier > 0) {
                    orders.push(actions.createMoveAction(terrains[pipe.first].terrainId, terrains[pipe.second].terrainId, 1));
                }
                if (playerTerrainsIndex.indexOf(pipe.second) !== -1 && playerTerrainsIndex.indexOf(pipe.first) === -1 && terrains[pipe.second].numberOfSoldier > 0) {
                    orders.push(actions.createMoveAction(terrains[pipe.second].terrainId, terrains[pipe.first].terrainId, 1));
                }
            }
            return orders;
        }
        return [];
    }
};