package bindings

import (
	"strings"
	"time"

	"github.com/relvacode/iso8601"
)

func MarshalDateTime(t *time.Time) ([]byte, error) {
	f := "\""+t.Format("2006-01-02T03:04:05-07:00")+"\""
	return []byte(f), nil
}

func UnmarshalDateTime(b []byte, t *time.Time) error {
	in := strings.ReplaceAll(string(b), "\"", "")
	time, err := iso8601.ParseString(in)
	*t = time
	return err
}
