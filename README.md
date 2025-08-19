# space-cli

Welcome to my cli-project for ReMarkable's technical process.

This project aims to explore API's that are out of this world.

Startup:

```sh
go build -o space-cli
```

Following the commands below you can get to know the latest stats in space exploration:

Get the upcoming launch stats (Data Sources: SpaceX):

```sh
./space-cli launches --upcoming
```

Get the launch stats between given dates (Data Sources: SpaceX):

```sh
./space-cli launches --start 2020-01-01 --end 2025-01-01
```

Get the launch stats for failed launches (Data Sources: SpaceX):

```sh
./space-cli launches --limit 5 --failed
```

Get the launch total costs (Data Sources: SpaceX):

```sh
./space-cli launches --cost
```

Get launches with near-earth asteroid data from NASA (Data Sources: SpaceX, NASA):

```sh
./space-cli launches --limit 5 --asteroids
```

Example combinations;

- Get the total cost of all failed launches between given dates (Data Sources: SpaceX):

```sh
./space-cli launches --start 2006-01-01 --end 2025-01-01 --failed --cost
```

- Get launch stats for the last 5 launches with location data (Data Sources: SpaceX):

```sh
./space-cli launches --limit 5 --launchpad
```

- Get launch stats for the last 5 launches with location and weather data (Data Sources: SpaceX, NASA):

```sh
./space-cli launches --limit 5 --launchpad --weather
```
