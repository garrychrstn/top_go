-- Composite index for SSO lookups, which are always scoped to an app:
-- app_id leads so WHERE app_id = $1 AND (username|email) = $2 can use it.
CREATE INDEX users_username_email_app_idx ON users (app_id, username, email);
