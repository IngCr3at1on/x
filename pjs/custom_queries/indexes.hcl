# The label is the name and usage message used for the cobra command.
query "indexes-for-table <table>" {
    // Description is the cobra.Command.Long and prints as part of the help text.
    description = "Returns a list of indexes for the provided table name."
    /* Use heredoc for multi-line strings. */
    template = <<EOF
SELECT * FROM pg_indexes
WHERE tablename = '{{ .Table }}';
EOF
    args = [ "Table" ]
}