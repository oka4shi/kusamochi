package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Post(url string, body string) (*http.Response, error) {
	// Replace " with \" and assemble json
	jsonData := fmt.Appendf([]byte{}, `{"content": "%s"}`, strings.ReplaceAll(strings.ReplaceAll(body, "\"", "\\\""), "\n", "\\n"))
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

type discordParams struct {
	Content     string       `json:"content"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
}

type multipart struct {
	ContentDisposition string
	ContentType        string
	Content            []byte
}

type File struct {
	Name        string
	Content     *bytes.Buffer
	Description string
}

func PostWithFiles(url string, body string, files []File) (*http.Response, error) {
	BOUNDARY := "KusamochiBoundary"

	attachments := []attachment{}
	for i, file := range files {
		attachments = append(attachments, attachment{
			ID:          i,
			Description: file.Description,
			Filename:    file.Name,
		})
	}

	parts := []multipart{}
	data := discordParams{
		Content:     body,
		Attachments: attachments,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	parts = append(parts, multipart{
		ContentDisposition: "form-data; name=\"payload_json\"",
		ContentType:        "application/json",
		Content:            jsonData,
	})

	for i, file := range files {
		content := file.Content.Bytes()
		parts = append(parts, multipart{
			ContentDisposition: fmt.Sprintf("form-data; name=\"files[%d]\"; filename=\"%s\"", i, file.Name),
			ContentType:        http.DetectContentType(content),
			Content:            content,
		})
	}

	reqBody := new(bytes.Buffer)
	for _, part := range parts {
		fmt.Fprintf(reqBody, "--%s\r\n", BOUNDARY)
		fmt.Fprintf(reqBody, "Content-Disposition: %s\r\n", part.ContentDisposition)
		fmt.Fprintf(reqBody, "Content-Type: %s\r\n", part.ContentType)
		reqBody.WriteString("\r\n")
		reqBody.Write(part.Content)
		reqBody.WriteString("\r\n")
	}

	reqBody.WriteString(fmt.Sprintf("--%s--\r\n", BOUNDARY))

	request, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", BOUNDARY))

	webhook := &http.Client{}
	response, err := webhook.Do(request)
	if err != nil {
		return nil, err
	}
	return response, err
}
