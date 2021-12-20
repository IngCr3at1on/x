# pjs

Yes as in pijamas; pjs started as an experiment with jackc/pgx to see if I could read data without knowing what it was ahead of time. In the end it became a somewhat useful CLI postgres client that auto-outputs to json.

When combined with jq I've found it to be very powerful.

## Install

    go install github.com/ingcr3at1on/x/pjs/cmd/pjs

## Use

See test-local.env for a list of recognized postgres env variables.
With appropriate variables set simple queries can be ran as follows:

    echo "SELECT * FROM pg_proc LIMIT 1;" | pjs | jq .

For more complex queries consider writing to a file and reading with cat.

    cat <file> | pjs | jq .

pjs also supports passing a dsn string directly instead of using env variables by use of the `--dsn` (or `-d`) flag.

## Tests

    docker-compose up -d
    sleep 3
    godotenv -f test-local.env go test -race ./...
    godotenv -f test-local-alt.env go test -race ./...
    docker-compose down
