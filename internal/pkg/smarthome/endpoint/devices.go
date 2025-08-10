package endpoint

import (
	"net/http"
	"net/url"
	"strconv"
)

func Devices(accessToken string, limit, offset int, ids ...string) *Endpoint {
	return &Endpoint{
		Method: "GET",
		URL:    MustURL("https://gateway.sbrdvc.xyz/gateway/v1/devices"),
		Headers: http.Header{
			"X-Auth-Jwt": {accessToken},
		},
		Query: url.Values{
			"pagination.limit":  []string{strconv.Itoa(limit)},
			"pagination.offset": []string{strconv.Itoa(offset)},
			"id":                ids,
		},
	}
}
