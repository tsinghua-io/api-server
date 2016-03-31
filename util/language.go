package util

import (
	"golang.org/x/text/language"
	"net/http"
)

var SupportedLanguageMatcher = language.NewMatcher([]language.Tag{
	language.SimplifiedChinese,
	language.English,
})

func Language(req *http.Request) language.Tag {
	tags, _, _ := language.ParseAcceptLanguage(req.Header.Get("Accept-Language"))
	tag, _, _ := SupportedLanguageMatcher.Match(tags...)
	return tag
}
