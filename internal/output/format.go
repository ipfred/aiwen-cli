package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/aiwen/aw-cli/errs"
)

const (
	FormatJSON   = "json"
	FormatNDJSON = "ndjson"
	FormatTable  = "table"
	FormatCSV    = "csv"
)

type Envelope struct {
	OK     bool   `json:"ok"`
	Action string `json:"action,omitempty"`
	IP     string `json:"ip,omitempty"`
	Data   any    `json:"data,omitempty"`
	Error  any    `json:"error,omitempty"`
}

func ValidateFormat(format string) error {
	switch format {
	case FormatJSON, FormatNDJSON, FormatTable, FormatCSV:
		return nil
	default:
		return errs.Validationf("unsupported format: %s", format)
	}
}

func Write(w io.Writer, value any, format string) error {
	if err := ValidateFormat(format); err != nil {
		return err
	}
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		return enc.Encode(value)
	case FormatNDJSON:
		return writeNDJSON(w, value)
	case FormatTable:
		return writeTable(w, value)
	case FormatCSV:
		return writeCSV(w, value)
	default:
		return nil
	}
}

func Filter(value any, expr string) any {
	if expr == "" || expr == "." {
		return value
	}
	path := strings.TrimPrefix(expr, ".")
	parts := strings.Split(path, ".")
	current := normalize(value)
	for _, part := range parts {
		obj, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = obj[part]
	}
	return current
}

func writeNDJSON(w io.Writer, value any) error {
	values, ok := value.([]any)
	if !ok {
		values = []any{value}
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for _, item := range values {
		if err := enc.Encode(item); err != nil {
			return err
		}
	}
	return nil
}

func writeTable(w io.Writer, value any) error {
	rows := flattenRows(value)
	if len(rows) == 0 {
		return nil
	}
	keys := sortedKeys(rows)
	fmt.Fprintln(w, strings.Join(keys, "\t"))
	for _, row := range rows {
		fields := make([]string, 0, len(keys))
		for _, key := range keys {
			fields = append(fields, fmt.Sprint(row[key]))
		}
		fmt.Fprintln(w, strings.Join(fields, "\t"))
	}
	return nil
}

func writeCSV(w io.Writer, value any) error {
	rows := flattenRows(value)
	if len(rows) == 0 {
		return nil
	}
	keys := sortedKeys(rows)
	writer := csv.NewWriter(w)
	if err := writer.Write(keys); err != nil {
		return err
	}
	for _, row := range rows {
		fields := make([]string, 0, len(keys))
		for _, key := range keys {
			fields = append(fields, fmt.Sprint(row[key]))
		}
		if err := writer.Write(fields); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func flattenRows(value any) []map[string]any {
	normalized := normalize(value)
	if values, ok := normalized.([]any); ok {
		rows := make([]map[string]any, 0, len(values))
		for _, item := range values {
			if row, ok := item.(map[string]any); ok {
				rows = append(rows, row)
			}
		}
		return rows
	}
	if row, ok := normalized.(map[string]any); ok {
		return []map[string]any{row}
	}
	return nil
}

func sortedKeys(rows []map[string]any) []string {
	seen := map[string]bool{}
	for _, row := range rows {
		for key := range row {
			seen[key] = true
		}
	}
	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func normalize(value any) any {
	data, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var out any
	if err := json.Unmarshal(data, &out); err != nil {
		return value
	}
	return out
}
