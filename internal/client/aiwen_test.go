package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiwen/aw-cli/internal/endpoint"
)

func TestBuildRequestLocIPv4(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key-123",
		Channel: "aw_cli",
	}
	req := QueryRequest{
		Action:   endpoint.ActionLoc,
		IP:       "8.8.8.8",
		Accuracy: "city",
	}
	httpReq, err := c.BuildRequest(context.Background(), req, false)
	if err != nil {
		t.Fatalf("BuildRequest error: %v", err)
	}
	if httpReq.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", httpReq.Method)
	}
	if httpReq.URL.Path != "/ip/geo/v1/city/" {
		t.Errorf("expected path /ip/geo/v1/city/, got %s", httpReq.URL.Path)
	}
	query := httpReq.URL.Query()
	if query.Get("ip") != "8.8.8.8" {
		t.Errorf("expected ip=8.8.8.8, got %s", query.Get("ip"))
	}
	if query.Get("key") != "test-key-123" {
		t.Errorf("expected key=test-key-123, got %s", query.Get("key"))
	}
	if query.Get("channel") != "aw_cli" {
		t.Errorf("expected channel=aw_cli, got %s", query.Get("channel"))
	}
}

func TestBuildRequestRedactsKey(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key-123",
		Channel: "aw_cli",
	}
	req := QueryRequest{
		Action: endpoint.ActionLoc,
		IP:     "8.8.8.8",
	}
	httpReq, err := c.BuildRequest(context.Background(), req, true)
	if err != nil {
		t.Fatalf("BuildRequest error: %v", err)
	}
	query := httpReq.URL.Query()
	key := query.Get("key")
	if key == "test-key-123" {
		t.Error("key should be redacted in preview mode")
	}
	if key != "tes***123" {
		t.Errorf("expected redacted key tes***123, got %s", key)
	}
}

func TestBuildRequestIPv4Only(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key",
		Channel: "aw_cli",
	}
	req := QueryRequest{
		Action: endpoint.ActionRisk,
		IP:     "2001:4860:4860::8888",
	}
	_, err := c.BuildRequest(context.Background(), req, false)
	if err == nil {
		t.Error("expected error for IPv6 on IPv4-only action")
	}
}

func TestBuildRequestWithCoordsys(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key",
		Channel: "aw_cli",
	}
	req := QueryRequest{
		Action:   endpoint.ActionLoc,
		IP:        "8.8.8.8",
		Accuracy: "city",
		CoordSys: "GCJ02",
	}
	httpReq, err := c.BuildRequest(context.Background(), req, false)
	if err != nil {
		t.Fatalf("BuildRequest error: %v", err)
	}
	query := httpReq.URL.Query()
	if query.Get("coordsys") != "GCJ02" {
		t.Errorf("expected coordsys=GCJ02, got %s", query.Get("coordsys"))
	}
}

func TestBuildRequestWithLang(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key",
		Channel: "aw_cli",
	}
	req := QueryRequest{
		Action: endpoint.ActionScene,
		IP:     "8.8.8.8",
		Lang:   "en",
	}
	httpReq, err := c.BuildRequest(context.Background(), req, false)
	if err != nil {
		t.Fatalf("BuildRequest error: %v", err)
	}
	query := httpReq.URL.Query()
	if query.Get("lang") != "en" {
		t.Errorf("expected lang=en, got %s", query.Get("lang"))
	}
}

func TestQueryMissingAPIKey(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "",
		Channel: "aw_cli",
	}
	_, err := c.Query(context.Background(), QueryRequest{
		Action: endpoint.ActionLoc,
		IP:     "8.8.8.8",
	})
	if err == nil {
		t.Error("expected error for missing API key")
	}
}

func TestQuerySuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"country":"US","city":"Mountain View"}`))
	}))
	defer server.Close()

	c := &AiwenClient{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Channel: "aw_cli",
		HTTP:    &http.Client{},
	}
	result, err := c.Query(context.Background(), QueryRequest{
		Action: endpoint.ActionLoc,
		IP:     "8.8.8.8",
	})
	if err != nil {
		t.Fatalf("Query error: %v", err)
	}
	if result.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", result.StatusCode)
	}
}

func TestQueryServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
	}))
	defer server.Close()

	c := &AiwenClient{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Channel: "aw_cli",
		HTTP:    &http.Client{},
	}
	_, err := c.Query(context.Background(), QueryRequest{
		Action: endpoint.ActionLoc,
		IP:     "8.8.8.8",
	})
	if err == nil {
		t.Error("expected error for server 500")
	}
}

func TestCurrentMissingAPIKey(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "",
		Channel: "aw_cli",
	}
	_, err := c.Current(context.Background(), QueryRequest{Action: endpoint.ActionCurrent})
	if err == nil {
		t.Error("expected error for missing API key")
	}
}

func TestPreview(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key-123",
		Channel: "aw_cli",
	}
	preview, err := c.Preview(context.Background(), QueryRequest{
		Action:   endpoint.ActionLoc,
		IP:       "8.8.8.8",
		Accuracy: "city",
	})
	if err != nil {
		t.Fatalf("Preview error: %v", err)
	}
	if preview.Method != "GET" {
		t.Errorf("expected method GET, got %s", preview.Method)
	}
	if preview.Query["ip"] != "8.8.8.8" {
		t.Errorf("expected ip=8.8.8.8, got %s", preview.Query["ip"])
	}
}

func TestInvalidIP(t *testing.T) {
	c := &AiwenClient{
		BaseURL: "https://api.ipplus360.com",
		APIKey:  "test-key",
	}
	_, err := c.BuildRequest(context.Background(), QueryRequest{
		Action: endpoint.ActionLoc,
		IP:     "not-an-ip",
	}, false)
	if err == nil {
		t.Error("expected error for invalid IP")
	}
}
