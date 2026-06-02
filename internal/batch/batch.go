package batch

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aiwen/aw-cli/internal/client"
	"github.com/aiwen/aw-cli/internal/iputil"
)

type Result struct {
	IP      string `json:"ip"`
	Action  string `json:"action"`
	OK      bool   `json:"ok"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	HTTPCode int   `json:"status_code,omitempty"`
}

type RunOptions struct {
	Concurrency int
	Retries     int
	DryRun      bool
}

func Run(ctx context.Context, c *client.AiwenClient, ips []string, actions []string, opts RunOptions) []Result {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 5
	}
	if opts.Retries < 0 {
		opts.Retries = 2
	}

	var jobs []job
	for _, ip := range ips {
		for _, action := range actions {
			jobs = append(jobs, job{IP: ip, Action: action})
		}
	}

	results := make([]Result, 0, len(jobs))
	ch := make(chan Result, len(jobs))
	jobsCh := make(chan job, len(jobs))

	var wg sync.WaitGroup
	for i := 0; i < opts.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsCh {
				ch <- executeJob(ctx, c, j, opts)
			}
		}()
	}

	for _, j := range jobs {
		jobsCh <- j
	}
	close(jobsCh)

	wg.Wait()
	close(ch)

	for r := range ch {
		results = append(results, r)
	}
	return results
}

func ParseFile(path, ipColumn string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if strings.HasSuffix(strings.ToLower(path), ".csv") {
		return parseCSV(f, ipColumn)
	}
	return parseText(f)
}

func FilterSpecial(ips []string) []string {
	filtered := make([]string, 0, len(ips))
	for _, ip := range ips {
		addr, err := iputil.Parse(ip)
		if err != nil {
			filtered = append(filtered, ip)
			continue
		}
		if !iputil.IsSpecial(addr.Addr) {
			filtered = append(filtered, ip)
		}
	}
	return filtered
}

type job struct {
	IP     string
	Action string
}

func executeJob(ctx context.Context, c *client.AiwenClient, j job, opts RunOptions) Result {
	req := client.QueryRequest{
		Action: j.Action,
		IP:     j.IP,
	}

	if opts.DryRun {
		preview, err := c.Preview(ctx, req)
		if err != nil {
			return Result{IP: j.IP, Action: j.Action, OK: false, Error: err.Error()}
		}
		return Result{IP: j.IP, Action: j.Action, OK: true, Data: preview}
	}

	var result *client.QueryResult
	var err error
	for attempt := 0; attempt <= opts.Retries; attempt++ {
		result, err = c.Query(ctx, req)
		if err == nil {
			break
		}
	}

	if err != nil {
		return Result{IP: j.IP, Action: j.Action, OK: false, Error: err.Error()}
	}

	var data any
	if err := json.Unmarshal(result.JSON, &data); err != nil {
		data = map[string]string{"raw": string(result.Raw)}
	}

	return Result{
		IP:       j.IP,
		Action:   j.Action,
		OK:       result.StatusCode < 400,
		Data:     data,
		HTTPCode: result.StatusCode,
	}
}

func parseText(f *os.File) ([]string, error) {
	var ips []string
	dec := json.NewDecoder(f)
	hasJSON := false

	for {
		var obj map[string]any
		if err := dec.Decode(&obj); err != nil {
			break
		}
		hasJSON = true
		if ip, ok := obj["ip"].(string); ok && ip != "" {
			ips = append(ips, ip)
		}
	}

	if hasJSON {
		return ips, nil
	}

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			ips = append(ips, line)
		}
	}
	return ips, nil
}

func parseCSV(f *os.File, ipColumn string) ([]string, error) {
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	if len(records) == 0 {
		return nil, nil
	}

	header := records[0]
	colIndex := -1
	for i, h := range header {
		if strings.EqualFold(strings.TrimSpace(h), ipColumn) {
			colIndex = i
			break
		}
	}
	if colIndex == -1 {
		if len(header) > 0 {
			colIndex = 0
		} else {
			return nil, fmt.Errorf("CSV has no columns")
		}
	}

	var ips []string
	for _, row := range records[1:] {
		if colIndex < len(row) {
			ip := strings.TrimSpace(row[colIndex])
			if ip != "" {
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}
