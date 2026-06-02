package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aiwen/aw-cli/errs"
	"github.com/aiwen/aw-cli/internal/core"
	"github.com/aiwen/aw-cli/internal/endpoint"
	"github.com/aiwen/aw-cli/internal/iputil"
)

const currentIPURL = "https://www.ipuu.net/ipuu/user/getIP"

type AiwenClient struct {
	BaseURL string
	APIKey  string
	Channel string
	HTTP    *http.Client
	ErrOut  io.Writer
}

type QueryRequest struct {
	Action   string `json:"action"`
	IP       string `json:"ip,omitempty"`
	Accuracy string `json:"accuracy,omitempty"`
	CoordSys string `json:"coordsys,omitempty"`
	Lang     string `json:"lang,omitempty"`
}

type QueryResult struct {
	StatusCode int             `json:"status_code"`
	Raw        []byte          `json:"-"`
	JSON       json.RawMessage `json:"json"`
}

type RequestPreview struct {
	Method string            `json:"method"`
	URL    string            `json:"url"`
	Query  map[string]string `json:"query"`
}

func (c *AiwenClient) Query(ctx context.Context, req QueryRequest) (*QueryResult, error) {
	if c.APIKey == "" {
		return nil, errs.Config("AIWEN_API_KEY is required")
	}
	httpReq, err := c.BuildRequest(ctx, req, false)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, errs.Network(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Network(err.Error())
	}
	if resp.StatusCode >= 500 {
		return nil, errs.Network("upstream server error")
	}
	var raw json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, errs.Parse("upstream returned non-JSON response")
	}
	result := &QueryResult{StatusCode: resp.StatusCode, Raw: body, JSON: raw}
	if resp.StatusCode >= 400 {
		return result, errs.API("upstream API returned an error")
	}
	return result, nil
}

func (c *AiwenClient) Current(ctx context.Context, req QueryRequest) (*QueryResult, error) {
	if c.APIKey == "" {
		return nil, errs.Config("AIWEN_API_KEY is required")
	}
	currentReq, err := http.NewRequestWithContext(ctx, http.MethodGet, currentIPURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient().Do(currentReq)
	if err != nil {
		return nil, errs.Network(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Network(err.Error())
	}
	var parsed struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, errs.Parse("current IP service returned non-JSON response")
	}
	if parsed.Data == "" {
		return nil, errs.API("current IP service did not return an IP")
	}
	req.Action = endpoint.ActionLoc
	req.IP = parsed.Data
	return c.Query(ctx, req)
}

func (c *AiwenClient) BuildRequest(ctx context.Context, req QueryRequest, redacted bool) (*http.Request, error) {
	if req.Action == endpoint.ActionCurrent {
		return http.NewRequestWithContext(ctx, http.MethodGet, currentIPURL, nil)
	}
	addr, err := iputil.Parse(req.IP)
	if err != nil {
		return nil, err
	}
	version := "ipv6"
	if addr.Addr.Is4() {
		version = "ipv4"
	}
	path, err := endpoint.Path(req.Action, version, req.Accuracy)
	if err != nil {
		return nil, err
	}
	base, err := url.Parse(strings.TrimRight(c.BaseURL, "/"))
	if err != nil {
		return nil, errs.Config("invalid base URL")
	}
	base.Path = strings.TrimRight(base.Path, "/") + path
	q := base.Query()
	q.Set("ip", req.IP)
	q.Set("channel", firstNonEmpty(c.Channel, core.DefaultChannel))
	if req.CoordSys != "" {
		q.Set("coordsys", req.CoordSys)
	}
	if req.Lang != "" {
		q.Set("lang", req.Lang)
	}
	if c.APIKey != "" {
		key := c.APIKey
		if redacted {
			key = core.RedactSecret(key)
		}
		q.Set("key", key)
	}
	base.RawQuery = q.Encode()
	return http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
}

func (c *AiwenClient) Preview(ctx context.Context, req QueryRequest) (RequestPreview, error) {
	httpReq, err := c.BuildRequest(ctx, req, true)
	if err != nil {
		return RequestPreview{}, err
	}
	query := map[string]string{}
	for key, values := range httpReq.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}
	return RequestPreview{
		Method: httpReq.Method,
		URL:    httpReq.URL.Scheme + "://" + httpReq.URL.Host + httpReq.URL.Path,
		Query:  query,
	}, nil
}

func (c *AiwenClient) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
