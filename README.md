# DB

DB is a toy key-value database written in Go.

## Usage

DB is a HTTP server that listens on port 4000.

### Set

To set a key-value pair, send a POST request to `http://localhost:4000/set?somekey=somevalue`.

### Get

To get the value of a key, send a GET request to `http://localhost:4000/get?key=somekey`.

## Development

### Requirements

* [Go](https://go.dev/dl/)
* [Just](https://github.com/casey/just#Installation)
* [Posting](https://github.com/darrenburns/posting#Installation)

### Running

`just dev` will start the server on port 4000.

`just test` will run the tests.

`just collection` will run the posting collection.

For more commands, run `just`.
