# ReMarkable-cli

This CLI project aims to explore API's that are out of this world.

Startup:

```sh
go build -o space-cli
```

Following the commands below you can get to know the latest stats in space exploration:

Get the upcoming launch stats from SpaceX:

```sh
./space-cli launches upcoming --limit 3
```

Get the past launch stats from SpaceX:

```sh
./space-cli launches past --limit 3
```
