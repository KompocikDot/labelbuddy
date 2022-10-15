package utils

import "regexp"

var TokenRe, _ 	= regexp.Compile(`"_token" value="(.{40})">`)
var LinksRe, _ = regexp.Compile(`window\.open\('https:\\\/\\\/restocks\.net\\\/en\\\/account\\\/sales\\\/send-label\\\/(.{7,8})'\)`)
var DatesRe, _ = regexp.Compile(`Ship before:\\n(\s*)(\d{2}\\/\d{2}\\/\d{2})`)
var LoginErrRe, _ = regexp.Compile(`combination is unknown`)