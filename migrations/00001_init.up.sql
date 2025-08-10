CREATE TABLE IF NOT EXISTS sessions (
    access_token                TEXT PRIMARY KEY,
    refresh_token               TEXT UNIQUE NOT NULL,

    otp_receiver_id             TEXT        NOT NULL,

    smart_home_access_token     TEXT        NOT NULL,
    smart_home_access_token_ttl INTERVAL    NOT NULL,

    smart_home_refresh_token    TEXT        NOT NULL,

    created_at                  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMP   NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS oauth_code (
    auth_code                      TEXT PRIMARY KEY,
    access_token                   TEXT REFERENCES sessions (access_token),
    smart_home_pkce_code_verifier  TEXT NOT NULL,
    smart_home_pkce_code_challenge TEXT NOT NULL,
    smart_home_otp_receiver_id     TEXT NOT NULL,
    smart_home_auth_operation_id   TEXT NOT NULL
);
