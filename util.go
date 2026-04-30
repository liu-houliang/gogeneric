package gogeneric

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// HTTPResponse wraps a decoded response body with HTTP metadata.
// Use res.Data for the body, res.StatusCode for the status, and res.Header for headers.
//
// Example:
//
//	res, err := HttpGetJson[MyResp](ctx, url, nil)
//	fmt.Println(res.Data, res.StatusCode, res.Header.Get("X-RateLimit-Remaining"))
type HTTPResponse[T any] struct {
	Data       T           // decoded JSON response body
	StatusCode int         // e.g. 200, 201, 204
	Header     http.Header // full response headers
}

// File represents one file in a multipart/form-data upload.
type File struct {
	FieldName string    // HTML form field name, e.g. "avatar"
	FileName  string    // file name sent in Content-Disposition, e.g. "photo.jpg"
	Reader    io.Reader // file content
}

// HTTPError is returned when the server responds with a non-2xx status code.
// It carries the status code and the raw response body, so callers can
// inspect or unmarshal structured error payloads returned by the API.
//
//	var httpErr *gogeneric.HTTPError
//	if errors.As(err, &httpErr) {
//	    var apiErr ApiErrorResp
//	    json.Unmarshal(httpErr.Body, &apiErr)
//	    fmt.Println(httpErr.StatusCode, apiErr.Message)
//	}
type HTTPError struct {
	StatusCode int         // e.g. 404, 500
	Status     string      // e.g. "404 Not Found"
	Header     http.Header // response headers (may contain rate-limit info etc.)
	Body       []byte      // raw response body (may be a JSON error payload)
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("gogeneric: HTTP %s", e.Status)
}

// ---------------------------------------------------------------------------
// Public HTTP functions
// ---------------------------------------------------------------------------

// HttpGetJson sends a GET request and decodes the JSON response into T.
// T may be a struct, slice, map, or pointer — any type json.Unmarshal supports.
//
//	res, err := HttpGetJson[[]User](ctx, "https://api.example.com/users", nil)
func HttpGetJson[T any](ctx context.Context, apiURL string, headers map[string]string) (*HTTPResponse[T], error) {
	return do[T](ctx, http.MethodGet, apiURL, "", nil, headers)
}

// HttpPostJson sends a POST with a JSON body and decodes the JSON response into T.
// request may be nil for an empty body.
//
//	res, err := HttpPostJson[User](ctx, url, CreateUserReq{Name: "alice"}, nil)
func HttpPostJson[T any](ctx context.Context, apiURL string, request any, headers map[string]string) (*HTTPResponse[T], error) {
	return doJSON[T](ctx, http.MethodPost, apiURL, request, headers)
}

// HttpPutJson sends a PUT request with a JSON body and decodes the JSON response into T.
//
//	res, err := HttpPutJson[User](ctx, url+"/1", UpdateReq{Name: "bob"}, nil)
func HttpPutJson[T any](ctx context.Context, apiURL string, request any, headers map[string]string) (*HTTPResponse[T], error) {
	return doJSON[T](ctx, http.MethodPut, apiURL, request, headers)
}

// HttpPatchJson sends a PATCH request with a JSON body and decodes the JSON response into T.
//
//	res, err := HttpPatchJson[User](ctx, url+"/1", PatchReq{Name: "carol"}, nil)
func HttpPatchJson[T any](ctx context.Context, apiURL string, request any, headers map[string]string) (*HTTPResponse[T], error) {
	return doJSON[T](ctx, http.MethodPatch, apiURL, request, headers)
}

// HttpDeleteJson sends a DELETE request and decodes the optional JSON response into T.
// If the server returns 204 No Content, Data will be the zero value of T.
//
//	res, err := HttpDeleteJson[DeletedResp](ctx, url+"/1", nil)
func HttpDeleteJson[T any](ctx context.Context, apiURL string, headers map[string]string) (*HTTPResponse[T], error) {
	return do[T](ctx, http.MethodDelete, apiURL, "", nil, headers)
}

// HttpPostForm sends a POST with application/x-www-form-urlencoded body
// and decodes the JSON response into T.
//
// Useful for OAuth token endpoints, legacy APIs, and HTML form submissions:
//
//	form := url.Values{"grant_type": {"password"}, "username": {"alice"}}
//	res, err := HttpPostForm[TokenResp](ctx, "https://auth.example.com/token", form, nil)
func HttpPostForm[T any](ctx context.Context, apiURL string, values url.Values, headers map[string]string) (*HTTPResponse[T], error) {
	return do[T](ctx, http.MethodPost, apiURL, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()), headers)
}

// HttpPostMultipart sends a POST with multipart/form-data body (file upload)
// and decodes the JSON response into T.
//
// files contains the files to upload; fields contains additional form text fields.
//
//	files := []gogeneric.File{{FieldName: "avatar", FileName: "photo.jpg", Reader: f}}
//	res, err := HttpPostMultipart[UploadResp](ctx, url, files, map[string]string{"userId": "1"}, nil)
func HttpPostMultipart[T any](ctx context.Context, apiURL string, files []File, fields map[string]string, headers map[string]string) (*HTTPResponse[T], error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	for k, v := range fields {
		if err := mw.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	for _, f := range files {
		part, err := mw.CreateFormFile(f.FieldName, f.FileName)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(part, f.Reader); err != nil {
			return nil, err
		}
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}

	return do[T](ctx, http.MethodPost, apiURL, mw.FormDataContentType(), &buf, headers)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// doJSON marshals request as JSON and delegates to do[T].
func doJSON[T any](ctx context.Context, method, apiURL string, request any, headers map[string]string) (*HTTPResponse[T], error) {
	var body io.Reader
	if request != nil {
		data, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}
	return do[T](ctx, method, apiURL, "application/json", body, headers)
}

// do is the single internal executor for all HTTP functions.
// It creates the request, sets headers, executes, checks status,
// reads the body, and unmarshals JSON into T.
func do[T any](ctx context.Context, method, apiURL, contentType string, body io.Reader, headers map[string]string) (*HTTPResponse[T], error) {
	req, err := http.NewRequestWithContext(ctx, method, apiURL, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read the body even on error so callers can inspect API error details.
		errBody, _ := io.ReadAll(resp.Body)
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Header:     resp.Header,
			Body:       errBody,
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result T
	if len(respBody) > 0 {
		if err = json.Unmarshal(respBody, &result); err != nil {
			return nil, err
		}
	}

	return &HTTPResponse[T]{
		Data:       result,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, nil
}
