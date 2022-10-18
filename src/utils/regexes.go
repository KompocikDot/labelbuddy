package utils

import "regexp"

var (
	TokenRe, _ 	= regexp.Compile(`"_token" value="(.{40})">`)
	LinksRe, _ = regexp.Compile(`window\.open\('https:\\\/\\\/restocks\.net\\\/en\\\/account\\\/sales\\\/send-label\\\/(.{7,8})'\)`)
	DatesRe, _ = regexp.Compile(`Ship before:\\n(\s*)(\d{2}\\/\d{2}\\/\d{2})`)
	LoginErrRe, _ = regexp.Compile(`combination is unknown`)
)