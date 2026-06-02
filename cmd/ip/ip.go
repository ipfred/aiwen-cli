package ip

import (
	"context"
	"encoding/json"

	"github.com/aiwen/aw-cli/errs"
	"github.com/aiwen/aw-cli/internal/client"
	"github.com/aiwen/aw-cli/internal/cmdutil"
	"github.com/aiwen/aw-cli/internal/core"
	"github.com/aiwen/aw-cli/internal/endpoint"
	"github.com/aiwen/aw-cli/internal/iputil"
	"github.com/aiwen/aw-cli/internal/output"
)

type IPQueryOptions struct {
	Factory  *cmdutil.Factory
	Ctx      context.Context
	Action   string
	IP       string
	Accuracy string
	CoordSys string
	Lang     string
	Format   string
	JqExpr   string
	DryRun   bool
}

func RunIPQuery(opts *IPQueryOptions) error {
	c, err := opts.Factory.Client()
	if err != nil {
		return err
	}

	req := client.QueryRequest{
		Action:   opts.Action,
		IP:       opts.IP,
		Accuracy: opts.Accuracy,
		CoordSys: opts.CoordSys,
		Lang:     opts.Lang,
	}

	if opts.DryRun {
		preview, err := c.Preview(opts.Ctx, req)
		if err != nil {
			return err
		}
		env := output.Envelope{
			OK:     true,
			Action: opts.Action,
			IP:     opts.IP,
			Data:   preview,
		}
		return output.Write(opts.Factory.IO.Out, env, opts.Format)
	}

	result, err := c.Query(opts.Ctx, req)
	if err != nil {
		return err
	}

	var data any
	if err := json.Unmarshal(result.JSON, &data); err != nil {
		return errs.Parse("failed to parse API response")
	}

	env := output.Envelope{
		OK:     true,
		Action: opts.Action,
		IP:     opts.IP,
		Data:   data,
	}

	filtered := output.Filter(env, opts.JqExpr)
	return output.Write(opts.Factory.IO.Out, filtered, opts.Format)
}

func RunCurrentQuery(opts *IPQueryOptions) error {
	c, err := opts.Factory.Client()
	if err != nil {
		return err
	}

	req := client.QueryRequest{
		Action: endpoint.ActionCurrent,
	}

	if opts.DryRun {
		preview := map[string]string{
			"method": "GET",
			"url":     "https://www.ipuu.net/ipuu/user/getIP",
			"note":    "resolves current IP first, then queries loc",
		}
		env := output.Envelope{
			OK:     true,
			Action: endpoint.ActionCurrent,
			Data:   preview,
		}
		return output.Write(opts.Factory.IO.Out, env, opts.Format)
	}

	result, err := c.Current(opts.Ctx, req)
	if err != nil {
		return err
	}

	var data any
	if err := json.Unmarshal(result.JSON, &data); err != nil {
		return errs.Parse("failed to parse API response")
	}

	env := output.Envelope{
		OK:     true,
		Action: endpoint.ActionCurrent,
		Data:   data,
	}

	filtered := output.Filter(env, opts.JqExpr)
	return output.Write(opts.Factory.IO.Out, filtered, opts.Format)
}

func ValidateIPForAction(ipStr, action string) error {
	addr, err := iputil.Parse(ipStr)
	if err != nil {
		return err
	}
	version := "ipv4"
	if addr.Addr.Is6() {
		version = "ipv6"
	}
	return endpoint.SupportsVersion(action, version)
}

func ResolveAccuracy(opts *IPQueryOptions) string {
	if opts.Accuracy != "" {
		return opts.Accuracy
	}
	cfg, err := opts.Factory.ResolveConfig()
	if err != nil {
		return core.DefaultAccuracy
	}
	addr, _ := iputil.Parse(opts.IP)
	if addr.Addr.Is4() {
		return cfg.IPv4Accuracy
	}
	return cfg.IPv6Accuracy
}
