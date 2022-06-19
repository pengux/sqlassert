BEGIN;

CREATE TABLE sqlassert_test(
	sku text NOT NULL,
	name text NOT NULL,
	CONSTRAINT sqlassert_test_pkey PRIMARY KEY (sku)
);

CREATE INDEX sqlassert_test_name_idx ON sqlassert_test(name);

INSERT INTO sqlassert_test(sku, name) VALUES ('sku1', 'name1');

COMMIT;
