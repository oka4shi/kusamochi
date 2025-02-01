package webhook

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func Post(url string, body string) (*http.Response, error) {
	// Replace " with \" and assemble json
	jsonData := []byte(fmt.Sprintf(`{"content": "%s"}`, strings.ReplaceAll(strings.ReplaceAll(body, "\"", "\\\""), "\n", "\\n")))
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	webhook := &http.Client{}
	response, err := webhook.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}
