package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"
)

type RequestType interface {
	SetGetTimeout(t time.Duration)
	SetPostTimeout(t time.Duration)
	Get(method string, params url.Values) (json.RawMessage, error)
	Post(method string, params url.Values, data map[string]NamedReader) (json.RawMessage, error)
}

const (
	// DefaultAPIURL is the default telegram API URL.
	DefaultAPIURL = "https://api.telegram.org"
	// DefaultGetTimeout is the default timeout to be set for a GET request.
	DefaultGetTimeout = time.Second * 3
	// DefaultPostTimeout is the default timeout to be set for a POST request.
	DefaultPostTimeout = time.Second * 10
)

type Response struct {
	// Ok: if true, request was successful, and result can be found in the Result field.
	// If false, error can be explained in the Description.
	Ok bool `json:"ok"`
	// Result: result of requests (if Ok)
	Result json.RawMessage `json:"result"`
	// ErrorCode: Integer error code of request. Subject to change in the future.
	ErrorCode int `json:"error_code"`
	// Description: contains a human readable description of the error result.
	Description string `json:"description"`
	// Parameters: Optional extra data which can help automatically handle the error.
	Parameters json.RawMessage `json:"parameters"`
}

type TelegramError struct {
	Method      string
	Params      url.Values
	Code        int
	Description string
}

func (t *TelegramError) Error() string {
	return fmt.Sprintf("unable to %s: %s", t.Method, t.Description)
}

type NamedReader interface {
	Name() string
	io.Reader
}

type NamedFile struct {
	File     io.Reader
	FileName string
}

func (nf NamedFile) Read(p []byte) (n int, err error) {
	return nf.File.Read(p)
}

func (nf NamedFile) Name() string {
	return nf.FileName
}

func fillBuffer(b *bytes.Buffer, data map[string]NamedReader) (string, error) {
	w := multipart.NewWriter(b)

	for field, file := range data {
		fileName := file.Name()
		if fileName == "" {
			fileName = field
		}

		part, err := w.CreateFormFile(field, fileName)
		if err != nil {
			return "", fmt.Errorf("failed to create form file for field %s and fileName %s: %w", field, fileName, err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return "", fmt.Errorf("failed to copy file contents of field %s to form: %w", field, err)
		}
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart form writer: %w", err)
	}

	return w.FormDataContentType(), nil
}
