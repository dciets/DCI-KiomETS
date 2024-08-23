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
 *
 * @param terrainFromId {string}
 * @param terrainToId {string}
 * @param quantity {number}
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
 *
 * @param terrainId {string}
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
 *
 * @param terrainId {string}
 */
function createBuildFactoryAction(terrainId) {
    const action = new Action();
    action.build = new BuildAction();
    action.build.terrainId = terrainId;
    action.build.terrainType = 1;
    action.actionType = 1;
    return action;
}

function createDemolishAction(terrainId) {
    const action = new Action();
    action.build = new BuildAction();
    action.build.terrainId = terrainId;
    action.build.terrainType = 2;
    action.actionType = 1;
    return action;
}

module.exports = {
    createMoveAction,
    createBuildBarricadeAction,
    createBuildFactoryAction,
    createDemolishAction
}