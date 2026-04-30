package gogeneric_test

// Tests for util.go
// Uses net/http/httptest — no live server required.
// Run: go test -run TestHttp .
// Run: go test -v -run TestHttp .

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/liu-houliang/gogeneric/v2"
)

type httpResp struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ctxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func newJSONServer(status int, body any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

// ---------------------------------------------------------------------------
// HttpGetJson
// ---------------------------------------------------------------------------

func TestHttpGetJson_Struct(t *testing.T) {
	srv := newJSONServer(200, httpResp{ID: 1, Name: "alice"})
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpGetJson[httpResp](ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.ID != 1 || res.Data.Name != "alice" {
		t.Fatalf("Data = %+v, want {1 alice}", res.Data)
	}
	if res.StatusCode != 200 {
		t.Fatalf("StatusCode = %d, want 200", res.StatusCode)
	}
}

// T can be a slice — impossible before the requireStruct check was removed.
func TestHttpGetJson_Slice(t *testing.T) {
	payload := []httpResp{{ID: 1, Name: "alice"}, {ID: 2, Name: "bob"}}
	srv := newJSONServer(200, payload)
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpGetJson[[]httpResp](ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("HttpGetJson[slice] error: %v", err)
	}
	if len(res.Data) != 2 {
		t.Fatalf("len = %d, want 2", len(res.Data))
	}
}

// T can be a map.
func TestHttpGetJson_Map(t *testing.T) {
	payload := map[string]any{"key": "value", "count": float64(3)}
	srv := newJSONServer(200, payload)
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpGetJson[map[string]any](ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("HttpGetJson[map] error: %v", err)
	}
	if res.Data["key"] != "value" {
		t.Fatalf("Data[key] = %v, want value", res.Data["key"])
	}
}

func TestHttpGetJson_ResponseHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "42")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: 1, Name: "alice"})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpGetJson[httpResp](ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := res.Header.Get("X-RateLimit-Remaining"); got != "42" {
		t.Fatalf("X-RateLimit-Remaining = %q, want 42", got)
	}
}

func TestHttpGetJson_WithHeaders(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	gogeneric.HttpGetJson[httpResp](ctx, srv.URL, map[string]string{"Authorization": "Bearer tok"})
	if gotAuth != "Bearer tok" {
		t.Fatalf("Authorization = %q, want Bearer tok", gotAuth)
	}
}

func TestHttpGetJson_NonOK(t *testing.T) {
	errPayload := `{"error":"not_found","message":"resource does not exist"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errPayload))
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	_, err := gogeneric.HttpGetJson[httpResp](ctx, srv.URL, nil)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}

	// error should be *HTTPError with body preserved
	var httpErr *gogeneric.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error type = %T, want *gogeneric.HTTPError", err)
	}
	if httpErr.StatusCode != 404 {
		t.Fatalf("HTTPError.StatusCode = %d, want 404", httpErr.StatusCode)
	}
	if string(httpErr.Body) != errPayload+"\n" && string(httpErr.Body) != errPayload {
		t.Fatalf("HTTPError.Body = %q, want JSON error payload", string(httpErr.Body))
	}
}

// ---------------------------------------------------------------------------
// HttpPostJson
// ---------------------------------------------------------------------------

type postReq struct {
	ID int `json:"id"`
}

func TestHttpPostJson_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		var req postReq
		json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: req.ID, Name: "echo"})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpPostJson[httpResp](ctx, srv.URL, postReq{ID: 42}, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.ID != 42 {
		t.Fatalf("ID = %d, want 42", res.Data.ID)
	}
}

func TestHttpPostJson_NilBody(t *testing.T) {
	srv := newJSONServer(200, httpResp{Name: "empty"})
	defer srv.Close()
	ctx, cancel := ctxTimeout()
	defer cancel()
	res, err := gogeneric.HttpPostJson[httpResp](ctx, srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.Name != "empty" {
		t.Fatalf("Name = %q, want empty", res.Data.Name)
	}
}

func TestHttpPostJson_NonOK(t *testing.T) {
	errPayload := `{"error":"unauthorized","message":"invalid token"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errPayload))
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	_, err := gogeneric.HttpPostJson[httpResp](ctx, srv.URL, nil, nil)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	var httpErr *gogeneric.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error type = %T, want *gogeneric.HTTPError", err)
	}
	if httpErr.StatusCode != 401 {
		t.Fatalf("HTTPError.StatusCode = %d, want 401", httpErr.StatusCode)
	}
	if len(httpErr.Body) == 0 {
		t.Fatal("HTTPError.Body should not be empty for 401 with JSON payload")
	}
}

// ---------------------------------------------------------------------------
// HttpPutJson / HttpPatchJson
// ---------------------------------------------------------------------------

func TestHttpPutJson_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: 1, Name: "updated"})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpPutJson[httpResp](ctx, srv.URL, postReq{ID: 1}, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.Name != "updated" {
		t.Fatalf("Name = %q, want updated", res.Data.Name)
	}
}

func TestHttpPatchJson_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Method = %s, want PATCH", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: 1, Name: "patched"})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpPatchJson[httpResp](ctx, srv.URL, map[string]string{"name": "patched"}, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.Name != "patched" {
		t.Fatalf("Name = %q, want patched", res.Data.Name)
	}
}

// ---------------------------------------------------------------------------
// HttpDeleteJson
// ---------------------------------------------------------------------------

func TestHttpDeleteJson_WithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %s, want DELETE", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: 1, Name: "deleted"})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpDeleteJson[httpResp](ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.Name != "deleted" {
		t.Fatalf("Name = %q, want deleted", res.Data.Name)
	}
}

func TestHttpDeleteJson_NoContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // 204, empty body
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	// T is struct but body is empty — should not error
	res, err := gogeneric.HttpDeleteJson[httpResp](ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("error on 204: %v", err)
	}
	if res.StatusCode != 204 {
		t.Fatalf("StatusCode = %d, want 204", res.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// HttpPostForm
// ---------------------------------------------------------------------------

func TestHttpPostForm_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q", ct)
		}
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: 1, Name: r.FormValue("name")})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	res, err := gogeneric.HttpPostForm[httpResp](ctx, srv.URL, url.Values{"name": {"alice"}}, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.Name != "alice" {
		t.Fatalf("Name = %q, want alice", res.Data.Name)
	}
}

// ---------------------------------------------------------------------------
// HttpPostMultipart
// ---------------------------------------------------------------------------

func TestHttpPostMultipart_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		userId := r.FormValue("userId")
		file, header, err := r.FormFile("avatar")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		content, _ := io.ReadAll(file)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{
			ID:   1,
			Name: userId + ":" + header.Filename + ":" + string(content),
		})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	files := []gogeneric.File{
		{FieldName: "avatar", FileName: "photo.jpg", Reader: strings.NewReader("image-bytes")},
	}
	res, err := gogeneric.HttpPostMultipart[httpResp](ctx, srv.URL, files, map[string]string{"userId": "7"}, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	want := "7:photo.jpg:image-bytes"
	if res.Data.Name != want {
		t.Fatalf("Name = %q, want %q", res.Data.Name, want)
	}
}

func TestHttpPostMultipart_MultipleFiles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		count := len(r.MultipartForm.File)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpResp{ID: count, Name: "ok"})
	}))
	defer srv.Close()

	ctx, cancel := ctxTimeout()
	defer cancel()

	files := []gogeneric.File{
		{FieldName: "f1", FileName: "a.txt", Reader: strings.NewReader("aaa")},
		{FieldName: "f2", FileName: "b.txt", Reader: strings.NewReader("bbb")},
	}
	res, err := gogeneric.HttpPostMultipart[httpResp](ctx, srv.URL, files, nil, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Data.ID != 2 {
		t.Fatalf("file count = %d, want 2", res.Data.ID)
	}
}

func TestHttpPostMultipart_NonOK(t *testing.T) {
	errPayload := `{"error":"server_error"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errPayload))
	}))
	defer srv.Close()
	ctx, cancel := ctxTimeout()
	defer cancel()
	_, err := gogeneric.HttpPostMultipart[httpResp](ctx, srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for 500")
	}
	var httpErr *gogeneric.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error type = %T, want *gogeneric.HTTPError", err)
	}
	if httpErr.StatusCode != 500 {
		t.Fatalf("HTTPError.StatusCode = %d, want 500", httpErr.StatusCode)
	}
}
