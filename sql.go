package sqlassert

const (
	postgresTableExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM information_schema.tables WHERE table_name = $1
);
`

	postgresColumnExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.columns
	WHERE
		table_name = $1
		AND column_name = $2
);
`
	postgresConstraintExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.constraint_column_usage
	WHERE
		table_name = $1
		AND constraint_name = $2
);
`

	postgresRowExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM "%s" WHERE %s
);
`

	postgresIndexExistsQuery = `
SELECT EXISTS(
	SELECT 1
	FROM
		pg_class t,
		pg_class i,
		pg_index ix
	WHERE
		t.oid = ix.indrelid
		AND i.oid = ix.indexrelid
		AND t.relkind = 'r'
		AND t.relname = $1
		AND i.relname = $2
);
`
)

const (
	mysqlTableExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM information_schema.tables WHERE table_name = ?
)
`
)
