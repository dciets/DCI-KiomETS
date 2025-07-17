import { createMoveAction } from './actions.js';
import { env } from './env.js';

export default class Agent {
    constructor() {

    }

    /**
     * Update function
     *
     * Players is a list of players.
     * Terrains is a list of terrains.
     * Pipes is a list of connection between terrains.
     * `pipe.first` and `pipe.second` are the index of both connected terrains.
     * `pipe.length` is the number of step for a soldier to pass through the link. A length of 4 takes 4 ticks to go through.
     * `pipe.soldier.length` is the position of a soldier group on the connection.
     * `terrain.terrainType` is either : 0 - barricade, 1 - factory, 2 - nothing or 3 - in construction
     *
     * @param data {
     *     {
     *         players: {name: string, color: string, numberOfKill: number, possessedTerrainsCount: number}[],
     *         terrains: {terrainId: string, terrainType: number, ownerIndex: number, numberOfSoldier: number}[],
     *         pipes: {length: number, first: number, second: number, soldiers: {ownerIndex: number, soldierCount: number, length: number}[]}[],
     *     }
     * }
     */
    update(data) {
        const playerName = env.playerName;

        const terrains = data.terrains;
        const players = data.players;
        const pipes = data.pipes;

        // Get the player
        const player = players.find(p => p.name === playerName);
        
        if (player) {
            // Get the terrains of the player
            const playerIndex = players.indexOf(player);
            const playerTerrains = terrains.filter(t => t.ownerIndex === playerIndex);

            // Get the indexes of the terrains
            const playerTerrainsIndex = playerTerrains.map(t => terrains.indexOf(t));

            const orders = [];
            for (const pipe of pipes) {
                // If the first terrain of a connection belongs to the player and have at least one soldier
                if (playerTerrainsIndex.indexOf(pipe.first) !== -1 && playerTerrainsIndex.indexOf(pipe.second) === -1 && terrains[pipe.first].numberOfSoldier > 0) {
                    // Send one soldier to the end of the connection
                    orders.push(createMoveAction(terrains[pipe.first].terrainId, terrains[pipe.second].terrainId, 1));
                }

                // If the second terrain of a connection belongs to the player and have at least one soldier
                if (playerTerrainsIndex.indexOf(pipe.second) !== -1 && playerTerrainsIndex.indexOf(pipe.first) === -1 && terrains[pipe.second].numberOfSoldier > 0) {
                    // Send one soldier to the end of the connection
                    orders.push(createMoveAction(terrains[pipe.second].terrainId, terrains[pipe.first].terrainId, 1));
                }
            }
            // Returns the order for each terrain of the player
            return orders;
        }

        // Returns no orders since the player was not found
        return [];
    }
};