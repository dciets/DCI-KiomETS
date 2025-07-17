/**
 * Environment variables for the agent
 */
class Env {
    _port = 10000;
    _magic = 0x11223344;

    _id = '';
    _player_name = '';

    _host = '';

    /**
     * Get the port number of the game server
     * @returns {number}
     */
    get port() {
        return this._port;
    }

    /**
     * Get the host address of the game server
     * @returns {string}
     */
    get host() {
        return this._host;
    }

    /**
     * Get the ID of the agent
     * @returns {string}
     */
    get id() {
        return this._id;
    }

    /**
     * Get the magic number for communication
     * @returns {number}
     */
    get magic() {
        return this._magic;
    }

    /**
     * Get the name of the agent
     * @returns {string}
     */
    get playerName() {
        return this._player_name;
    }

    set id(value) {
        this._id = value;
    }

    set playerName(value) {
        this._player_name = value;
    }
}

export const env = new Env();