package utils

import (
	"net/http"
	"net/http/cookiejar"
)

var (
	jar, _ = cookiejar.New(nil)
	HttpClient = &http.Client{
		Jar: jar,
	}
)