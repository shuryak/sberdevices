package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type Endpoint struct {
	Method  string
	Headers http.Header
	URL     *url.URL
	Query   url.Values
	Body    []byte
}

func (e *Endpoint) PreparedURL() string {
	var result url.URL
	result = *e.URL
	result.RawQuery = e.Query.Encode()
	return result.String()
}

func (e *Endpoint) CURL(cookieJar http.CookieJar, headers http.Header) string {
	if e.URL == nil {
		return ""
	}

	parts := []string{"curl"}

	if e.URL.Scheme == "https" {
		parts = append(parts, "-k")
	}

	parts = append(parts, "-X", e.Method)

	if e.Body != nil {
		if len(e.Body) != 0 {
			parts = append(parts, "-d", bashEscape(string(e.Body)))
		}
	}

	var headersKeys []string
	for key := range headers {
		headersKeys = append(headersKeys, key)
	}

	sort.Strings(headersKeys)

	preparedCookies := ""

	// TODO: additional cookie parameters
	if cookieJar != nil {
		cookies := cookieJar.Cookies(e.URL)

		for i := 0; i < len(cookies); i++ {
			if cookies[i].Value == "" {
				continue
			}

			preparedCookies += fmt.Sprintf("%s=%s", cookies[i].Name, cookies[i].Value)
			if i != len(cookies)-1 {
				preparedCookies += "; "
			}
		}

		if customCookies := headers.Get("Cookie"); customCookies != "" {
			preparedCookies += fmt.Sprintf("; %s", customCookies)
		}

		if headers.Get("Cookie") == "" && len(preparedCookies) > 0 {
			headersKeys = append(headersKeys, "Cookie")
		}
	}

	for _, key := range headersKeys {
		value := ""

		if key == "Cookie" {
			value = preparedCookies
		} else {
			value = strings.Join((headers)[key], " ")
		}

		parts = append(parts, "-H", bashEscape(fmt.Sprintf("%s: %s", key, value)))
	}

	parts = append(parts, bashEscape(e.URL.String()))

	return strings.Join(parts, " ")
}

func MustURL(rawURL string) *url.URL {
	u, _ := url.Parse(rawURL)
	return u
}

func MustJSON(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}

func bashEscape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}
