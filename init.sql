DROP TABLE IF EXISTS tax_allowance;
DROP TABLE IF EXISTS tax_level;

CREATE TABLE tax_allowance
(
    id                   bigserial NOT NULL,
    created_at           timestamptz NULL,
    updated_at           timestamptz NULL,
    deleted_at           timestamptz NULL,
    allowance_type       text      NOT NULL,
    min_allowance_amount numeric   NOT NULL,
    max_allowance_amount numeric   NOT NULL,
    CONSTRAINT tax_allowance_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_tax_allowance_deleted_at ON public.tax_allowance USING btree (deleted_at);

CREATE TABLE tax_level
(
    id          bigserial      NOT NULL,
    created_at  timestamptz NULL,
    updated_at  timestamptz NULL,
    deleted_at  timestamptz NULL,
    min_income  numeric(10, 2) NOT NULL,
    max_income  numeric(10, 2) NOT NULL,
    tax_percent numeric(10, 2) NOT NULL,
    CONSTRAINT tax_level_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_tax_level_deleted_at ON public.tax_level USING btree (deleted_at);

INSERT INTO tax_allowance (allowance_type, min_allowance_amount, max_allowance_amount)
VALUES ('personal', 60000.00, 60000.00),
       ('donation', 0.00, 100000.00),
       ('k-receipt', 0.00, 50000.00);

INSERT INTO tax_level (min_income, max_income, tax_percent)
VALUES (0.00, 150000.00, 0.00),
       (150001.00, 500000.00, 10.00),
       (500001.00, 1000000.00, 15.00),
       (1000001.00, 2000000.00, 20.00),
       (2000001.00, 2000001.00, 35.00)