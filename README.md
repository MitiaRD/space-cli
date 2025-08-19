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
```

Get the launch stats between given dates:

```sh
./space-cli launches --start 2020-01-01 --end 2025-01-01
```

Get the launch stats for failed launches:

```sh
./space-cli launches --failed
```

Get the launch total costs:

```sh
./space-cli launches --cost
```

Example combinations;

- Get the total cost of all failed launches between given dates

```sh
./space-cli launches --start 2006-01-01 --end 2025-01-01 --failed --cost
```

- Get launch stats for the last 5 launches with location data

```sh
./space-cli launches --limit 5 --launchpad
```

- Get launch stats for the last 5 launches with location and weather data from Nasa

```sh
./space-cli launches --limit 5 --launchpad --weather
```
