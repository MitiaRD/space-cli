# ReMarkable-cli

This CLI project aims to explore API's that are out of this world.

Startup:

```sh
go build -o space-cli
```

Following the commands below you can get to know the latest stats in space exploration:

Get the upcoming launch stats from SpaceX:

```sh
./space-cli launches --upcoming

./space-cli launches --start 2006-01-01 --end 2025-01-01

./space-cli launches --failed

./space-cli launches --cost

./space-cli launches --start 2006-01-01 --end 2025-01-01 --failed --cost
```
