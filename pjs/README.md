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

## Custom Queries

Experimental support for custom queries is baked into the cobra functionality; because of this the location for storing custom queries cannot be controlled from flags. The env value `PJS_CUSTOM_QUERIES` can be set to use a specific directory otherwise the default is `$HOME/.pjs/queries`.

Custom queries are read via HCL. The `template` parameter is a Go template, the template values correspond to the args which are processed in order against the input of the command.

A simple query is included in `custom_queries`. A more complex query might use heredoc to write the template parameter across multiple lines.

    # The label is the name and usage message used for the cobra command.
    query "indexes-for-table <table>" {
        // Description is the cobra.Command.Long and prints as part of the help text.
        description = "Returns a list of indexes for the provided table name."
        template = "SELECT * FROM pg_indexes WHERE tablename = '{{ .Table }}';"
        args = [ "Table" ]
    }

A more complex query would probably use heredoc to write the template parameter across multiple lines.

    query "indexes-for-table <table>" {
        description = "Returns a list of indexes for the provided table name."
        /* Use heredoc for multi-line strings. */
        template = <<EOF
    SELECT * FROM pg_indexes
    WHERE tablename = '{{ .Table }}';
    EOF
        args = [ "Table" ]
    }

See [custom_queries](https://github.com/IngCr3at1on/x/tree/master/pjs/custom_queries)

## Tests

    docker-compose up -d
    sleep 3
    godotenv -f test-local.env go test -race ./...
    godotenv -f test-local-alt.env go test -race ./...
    docker-compose down
