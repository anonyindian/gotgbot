package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BuiltinHttp struct {
	ApiUrl      string
	BotToken    string
	Client      http.Client
	PostTimeout time.Duration
	GetTimeout  time.Duration
}

func (h *BuiltinHttp) Get(method string, params url.Values) (json.RawMessage, error) {
	if h.GetTimeout == 0 {
		h.GetTimeout = DefaultGetTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.GetTimeout)
	defer cancel()

	return h.GetWithContext(ctx, method, params)
}

// GetWithContext allows sending a Get request with an existing context.
func (h *BuiltinHttp) GetWithContext(ctx context.Context, method string, params url.Values) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.methodEnpoint(method), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build GET request to %s: %w", method, err)
	}

	req.URL.RawQuery = params.Encode()

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GET request to %s: %w", method, err)
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode GET request to %s: %w", method, err)
	}

	if !r.Ok {
		return nil, &TelegramError{
			Method:      method,
			Params:      params,
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	return r.Result, nil
}

func (h *BuiltinHttp) Post(method string, params url.Values, data map[string]NamedReader) (json.RawMessage, error) {
	if h.PostTimeout == 0 {
		h.PostTimeout = DefaultPostTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.PostTimeout)
	defer cancel()

	return h.PostWithContext(ctx, method, params, data)
}

// PostWithContext allows sending a Post request with an existing context.
func (h *BuiltinHttp) PostWithContext(ctx context.Context, method string, params url.Values, data map[string]NamedReader) (json.RawMessage, error) {
	b := &bytes.Buffer{}
	contentType := "application/json"

	if len(data) > 0 {
		var err error
		contentType, err = fillBuffer(b, data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.methodEnpoint(method), b)
	if err != nil {
		return nil, fmt.Errorf("failed to build POST request to %s: %w", method, err)
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", contentType)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute POST request to %s: %w", method, err)
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode POST request to %s: %w", method, err)
	}

	if !r.Ok {
		return nil, &TelegramError{
			Method:      method,
			Params:      params,
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	return r.Result, nil
}

// GetAPIURL returns the currently used API endpoint.
func (h *BuiltinHttp) GetAPIURL() string {
	if h.ApiUrl == "" {
		return DefaultAPIURL
	}
	// Trim suffix to ensure consistent output
	return strings.TrimSuffix(h.ApiUrl, "/")
}

func (h *BuiltinHttp) methodEnpoint(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", h.GetAPIURL(), h.BotToken, method)
}

func (h *BuiltinHttp) SetGetTimeout(t time.Duration) {
	h.GetTimeout = t
}

func (h *BuiltinHttp) SetPostTimeout(t time.Duration) {
	h.GetTimeout = t
}
