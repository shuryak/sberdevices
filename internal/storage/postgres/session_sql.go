package postgres

const (
	sessionGetByAccessTokenSQL = `SELECT access_token, refresh_token, otp_receiver_id, smart_home_access_token, 
smart_home_access_token_ttl, smart_home_refresh_token, created_at, updated_at FROM sessions WHERE access_token = $1 
LIMIT 1;`

	sessionGetByRefreshTokenSQL = `SELECT access_token, refresh_token, otp_receiver_id, smart_home_access_token, 
smart_home_access_token_ttl, smart_home_refresh_token, created_at, updated_at FROM sessions WHERE access_token = $1 
LIMIT 1;`

	// $1 – old access token, $2 – new access token
	sessionUpsertSQL = `INSERT INTO sessions (access_token, refresh_token, otp_receiver_id, smart_home_access_token, 
smart_home_access_token_ttl, smart_home_refresh_token, created_at, updated_at) VALUES ($1, $3, $4, $5, $6, $7, $8, $9) 
ON CONFLICT (access_token) DO UPDATE SET access_token = $2, refresh_token = $3, smart_home_access_token = $5, 
smart_home_access_token_ttl = $6, smart_home_refresh_token = $7, updated_at = $9 RETURNING access_token, refresh_token, 
otp_receiver_id, smart_home_access_token, smart_home_access_token_ttl, smart_home_refresh_token, created_at, 
updated_at;`
)
