# SingularityNET Reputation Adapater
### Multi Party Escrow

## Development

These instructions are intended to facilitate the development and testing of SingularityNET Reputation Adapter. Users interested in
deploying the adapter should install the appropriate binary as
[released](#release).

### Prerequisites

* [Go 1.10+](https://golang.org/dl/)
* [postgresql](https://www.postgresql.org/download/)
* [Dep 0.4.1+](https://github.com/golang/dep#installation)
* [Node 8+ w/npm](https://nodejs.org/en/download/)
* [Yarn 1.12+](https://yarnpkg.com/lang/en/docs/install/)
* [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/downloads)


### Installing

* Clone the git repository
```bash
$ git clone git@github.com:singnet/reputation-adapter.git
$ cd reputation-adapter
```

* Install development/test dependencies
```bash
$ ./scripts/install
```

* Build reputation-adapter (on Linux)
```bash
$ ./scripts/build linux amd64
```

* Build reputation-adapter (on MacOSX)
```bash
$ ./scripts/build darwin amd64
```

* Build reputation-adapter for Linux (with Docker)
```bash
$ ./scripts/buildlinux
```

#### Run Deamon 

```bash
$ ./build/reputation-adapter-linux-amd64
```


## Run Postgres 

Run a `postgresql` instance and run [this SQL query](https://github.com/singnet/reputation-adapter/blob/master/resources/postgres/create_table.sql) against 

