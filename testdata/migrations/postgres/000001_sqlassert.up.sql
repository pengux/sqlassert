CREATE TABLE IF NOT EXISTS sqlassert_test(
	sku text UNIQUE NOT NULL,
	name text NOT NULL,
	CONSTRAINT sqlassert_test_pkey PRIMARY KEY (sku)
);

CREATE INDEX IF NOT EXISTS sqlassert_test_name_idx ON sqlassert_test(name);

INSERT INTO sqlassert_test(sku, name) VALUES ('sku1', 'name1');

