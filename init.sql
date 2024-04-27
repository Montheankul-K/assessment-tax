INSERT INTO tax_allowance (allowance_type, min_allowance_amount, max_allowance_amount)
VALUES ('personal', 60000.00, 100000.00), ('donation', 0.00, 100000.00), ('k-receipt', 0.00, 50000.00);

INSERT INTO tax_level (min_income, max_income, tax_percent)
VALUES (0.00, 150000.00, 0.00), (150001.00, 500000.00, 10.00), (500001.00, 1000000.00, 15.00), (1000001.00, 2000000.00, 20.00), (2000001.00, 2000001.00, 35.00)