BEGIN;

CREATE TABLE sqlassert_test (
    sku varchar(255) PRIMARY KEY,
    name varchar(255)
);

CREATE INDEX sqlassert_test_name_idx ON sqlassert_test(name);

INSERT INTO sqlassert_test(sku, name) VALUES ('sku1', 'name1');

COMMIT;
