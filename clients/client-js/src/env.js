class Env {

    _port = 10000;
    _host = '127.0.0.1';
    _id = 'de1cc5ba-56c9-11f0-bb34-3a457b3ca9d3';
    _magic = 0x11223344;
    _player_name = 'test';

    get port() {
        return this._port;
    }

    get host() {
        return this._host;
    }

    get id() {
        return this._id;
    }

    get magic() {
        return this._magic;
    }

    get playerName() {
        return this._player_name;
    }

    set port(value) {
        this._port = value;
    }

    set host(value) {
        this._host = value;
    }

    set id(value) {
        this._id = value;
    }

    set magic(value) {
        this._magic = value;
    }

    set playerName(value) {
        this._player_name = value;
    }
}

export const env = new Env();