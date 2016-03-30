package util

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func TrimmedTexts(s *goquery.Selection) []string {
	return s.Map(func(_ int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	})
}
