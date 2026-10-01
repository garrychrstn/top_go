CREATE OR REPLACE PROCEDURE sp_create_customer(
    IN p_name TEXT,
    IN p_phone_number TEXT,
    IN p_address TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO customers (name, phone_number, address)
    VALUES (p_name, p_phone_number, p_address);
END;
$$;
