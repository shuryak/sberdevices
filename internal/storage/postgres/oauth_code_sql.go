package postgres

const (
	oauthCodeCreateSQL = `INSERT INTO oauth_code (auth_code, smart_home_pkce_code_verifier, 
smart_home_pkce_code_challenge, smart_home_otp_receiver_id, smart_home_auth_operation_id) VALUES ($1, $2, $3, $4, $5);`

	oauthCodeGetSQL = `SELECT auth_code, access_token, smart_home_pkce_code_verifier, smart_home_pkce_code_challenge,
smart_home_otp_receiver_id, smart_home_auth_operation_id FROM oauth_code WHERE auth_code = $1;`

	oauthCodeSetAccessTokenSQL = `UPDATE oauth_code SET access_token = $1 WHERE auth_code = $2;`

	oauthCodeDeleteRowAndGetSessionSQL = `DELETE FROM oauth_code USING sessions WHERE oauth_code.auth_code = $1 AND 
oauth_code.access_token = sessions.access_token RETURNING sessions.access_token, sessions.refresh_token, 
sessions.otp_receiver_id, sessions.smart_home_access_token, sessions.smart_home_access_token_ttl, 
sessions.smart_home_refresh_token, sessions.created_at, sessions.updated_at;`
)
