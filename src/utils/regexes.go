package utils

import "regexp"

var (
	TokenRe, _ 	= regexp.Compile(`"_token" value="(.{40})">`)
	LinksRe, _ = regexp.Compile(`window\.open\('https:\\\/\\\/restocks\.net\\\/en\\\/account\\\/sales\\\/send-label\\\/(.{6,8})'\)`)
	DatesRe, _ = regexp.Compile(`Ship before:\\n(\s*)(\d{2}\\/\d{2}\\/\d{2})`)
	SoldItemsDatesRe, _ = regexp.Compile(`<td>\\n(\s*)(\d{2}\\/\d{2}\\/\d{2})`)
	ItemIdRe, _ = regexp.Compile(`ID: (\d{6,8})\\n`)
	LoginErrRe, _ = regexp.Compile(`combination is unknown`)
	EndOfPagesRe, _ = regexp.Compile(`no__listings__notice`)
	PriceRe, _ = regexp.Compile(`(\\u20ac|z\\u0142|\\u0024) (\d{2,}|\d{1,2}.\d{3,})\\n`) 
	// lookup for €, ł and $ signs and amount of money
	ConsignPriceRe, _ = regexp.Compile(`(\\u20ac|z\\u0142|\\u0024) (\d{2,}|\d{1,2}.\d{3,})`)
	ItemNameRe, _ = regexp.Compile(`<span>(.*?)<\\/span>`)
	SizeRe, _ = regexp.Compile(`EU: (\d+ \\u00bd|\d+ \\u2154|d+ \\u2153|\d+)`) // lookup for integers, ½, ⅔ and ⅓ signs
)