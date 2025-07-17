class MoveAction {
    /**
     * @type {string}
     */
    fromId
    /**
     * @type {string}
     */
    toId
    /**
     * @type {number}
     */
    quantity

    serialize() {
        return {fromId: this.fromId, toId: this.toId, quantity: this.quantity};
    }
}

class BuildAction {
    /**
     * @type {string}
     */
    terrainId
    /**
     * @type {number}
     */
    terrainType

    serialize() {
        return {terrainId: this.terrainId, terrainType: this.terrainType};
    }
}

class Action {
    /**
     * @type {number}
     */
    actionType
    /**
     * @type {MoveAction|null}
     */
    move
    /**
     * @type {BuildAction|null}
     */
    build
    serialise() {
        return {
            actionType: this.actionType,
            move: this.move?.serialize(),
            build: this.build?.serialize(),
        };
    }
}

/**
 * Create a move action to send to the server
 * @param terrainFromId {string} Source terrain id
 * @param terrainToId {string} Target adjacent terrain id
 * @param quantity {number} Quantity of soldier to send
 * @returns {Action}
 */
function createMoveAction(terrainFromId, terrainToId, quantity) {
    const action = new Action();
    action.move = new MoveAction();
    action.move.fromId = terrainFromId;
    action.move.toId = terrainToId;
    action.move.quantity = quantity;
    action.actionType = 0;
    return action;
}

/**
 * Create a build barricade action to send to the server
 * @param terrainId {string} Terrain id to build the barricade
 * @returns {Action}
 */
function createBuildBarricadeAction(terrainId) {
    const action = new Action();
    action.build = new BuildAction();
    action.build.terrainId = terrainId;
    action.build.terrainType = 0;
    action.actionType = 1;
    return action;
}

/**
 * Create a build factory action to send to the server
 * @param terrainId {string} Terrain id to build the factory
 * @returns {Action}
 */
function createBuildFactoryAction(terrainId) {
    const action = new Action();
    action.build = new BuildAction();
    action.build.terrainId = terrainId;
    action.build.terrainType = 1;
    action.actionType = 1;
    return action;
}

/**
 * Create a demolish action to send to the server
 * @param terrainId Terrain id to destroy a building
 * @returns {Action}
 */
function createDemolishAction(terrainId) {
    const action = new Action();
    action.build = new BuildAction();
    action.build.terrainId = terrainId;
    action.build.terrainType = 2;
    action.actionType = 1;
    return action;
}

export {
    createMoveAction,
    createBuildBarricadeAction,
    createBuildFactoryAction,
    createDemolishAction
}