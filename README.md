# timescale-coding-challenge
This is a command line tool that can be used to benchmark `SELECT` query performance across multiple workers against a
TimescaleDB instance. The tool takes as its input either a CSV-formatted file or standard input and a flag to specify
the number of concurrent workers. It starts processing queries as soon as each line from the CSV input is consumed.
After processing all the queries, the tool outputs a summary with the following statistics:
- Number of queries
- Total query time
- Minimum query time
- Maximum query time
- Average query time
- Median query time

## Input
The tool accepts input in the form of a CSV file or standard input. The input must be in CSV (comma-separated values)
format with the following columns:
- hostname (`host_000001`)
- start_time (`2017-01-01 08:59:22`)
- end_time (`2017-01-01 09:59:22`)

## Prerequisites
- [ ] Install [Git](https://git-scm.com/downloads) + [Go](https://go.dev/doc/install)
- [ ] Install [Docker Engine](https://docs.docker.com/engine/install)
- [ ] Create `.env` file in the root directory, using the [.env.example](.env.example) file as a template
- [ ] Prepare a CSV file with input data in the format specified [here](#input)

## Usage
### Docker Compose
To run the tool using Docker Compose, execute the following command:
```
make start
```
This command will perform the following steps in detached mode:
- Build a Docker image of the application
- Start the TimescaleDB container
- Perform database migrations
- Run the application with the configuration specified in the `.env` file and the input data from the CSV source

After the application has finished processing the input data, the results can be printed by executing the following command:
```
make print-results
```

To stop the database and remove all associated containers, execute the following command:
```
make stop
```

### Local
To use the tool locally, execute the following command:
```
make start-db
```
This command will perform the following steps in detached mode:
- Start the TimescaleDB container
- Perform database migrations

After starting the database, you can build the tool locally by executing the following command:
```
make build
```
This command will build the application binary in the `build` directory.

In order for the application to connect to the database container, you need to override the application configuration.
This can be done by creating a `.env.dev` file in the root directory with the following content:
```
POSTGRES_HOST=localhost
```
The application will merge the configuration from the `.env` file with the configuration from the `.env.dev` file.

After building the binary, you can perform the query benchmarks by executing the following command:
```
# Replace the flags with the desired values
ENV=DEV ./build/ts benchmark --workers=3 --config=query_params.csv
```
This command will perform the query benchmarks with 3 workers and the config specified in the `query_params.csv` file.
Finally, it will output the sorted benchmark results and query statistics in table format:
```
┏━━━━━┳━━━━━━━━━━━━━┓
┃     ┃    DURATION ┃
┣━━━━━╋━━━━━━━━━━━━━┫
┃   1 ┃   811.042µs ┃
┣━━━━━╋━━━━━━━━━━━━━┫
┃   N ┃         ... ┃
┣━━━━━╋━━━━━━━━━━━━━┫
┃ 200 ┃ 27.842125ms ┃
┗━━━━━┻━━━━━━━━━━━━━┛
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃       QUERY STATISTICS      ┃
┣━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━┫
┃ Count        ┃          200 ┃
┣━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━┫
┃ Total Time   ┃ 258.889377ms ┃
┣━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━┫
┃ Minimum Time ┃    811.042µs ┃
┣━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━┫
┃ Maximum Time ┃  27.842125ms ┃
┣━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━┫
┃ Average Time ┃   1.294446ms ┃
┣━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━┫
┃ Median Time  ┃    896.229µs ┃
┗━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━┛
```

To stop the database and remove all associated containers, execute the following command:
```
make stop
```

