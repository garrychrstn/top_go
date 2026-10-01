CREATE OR REPLACE PROCEDURE sp_create_user(
    IN p_username TEXT,
    IN p_password TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO users (username, password)
    VALUES (p_username, p_password);
END;
$$;
