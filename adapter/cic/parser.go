package cic

import (
	"io"
	"time"
)

type parser interface {
	parse(reader io.Reader, info interface{}, langCode string) error
}

func parseRegDate(regDate int64) string {
	return time.Unix(regDate/1000, 0).Format("2006-01-02T15:04:05+0800")
}
