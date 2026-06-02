package cmd

import (
	"context"
	"io"
	"os"

	"github.com/aiwen/aw-cli/errs"
	"github.com/aiwen/aw-cli/internal/batch"
	"github.com/aiwen/aw-cli/internal/cmdutil"
	"github.com/aiwen/aw-cli/internal/endpoint"
	"github.com/aiwen/aw-cli/internal/output"
	"github.com/spf13/cobra"
)

type BatchOptions struct {
	Factory        *cmdutil.Factory
	Action         string
	IPColumn       string
	OutputPath     string
	Concurrency    int
	Retries        int
	IncludePrivate bool
	InputFile      string
	Format         string
	JqExpr         string
	DryRun         bool
}

func NewCmdBatch(f *cmdutil.Factory) *cobra.Command {
	opts := &BatchOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "batch <file>",
		Short: "Batch query IPs from a file (txt, csv, or jsonl)",
		Long:  "Read IPs from a text, CSV, or JSONL file and query them concurrently. Supports txt (one IP per line), CSV (with --ip-column), and JSONL formats.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.InputFile = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			return runBatch(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Action, "action", "loc", "query action: loc, scene, whois, asn, host, risk, identity, industry, all")
	cmd.Flags().StringVar(&opts.IPColumn, "ip-column", "ip", "column name for IP addresses in CSV input")
	cmd.Flags().StringVarP(&opts.OutputPath, "output", "o", "", "output file path (default: stdout)")
	cmd.Flags().IntVar(&opts.Concurrency, "concurrency", 5, "number of concurrent requests")
	cmd.Flags().IntVar(&opts.Retries, "retries", 2, "number of retries for network errors")
	cmd.Flags().BoolVar(&opts.IncludePrivate, "include-private", false, "include private/reserved IP addresses")

	return cmd
}

func runBatch(opts *BatchOptions) error {
	actions := resolveActions(opts.Action)

	ips, err := batch.ParseFile(opts.InputFile, opts.IPColumn)
	if err != nil {
		return err
	}

	if !opts.IncludePrivate {
		ips = batch.FilterSpecial(ips)
	}

	c, err := opts.Factory.Client()
	if err != nil {
		return err
	}

	var out io.Writer = opts.Factory.IO.Out
	if opts.OutputPath != "" {
		f, err := os.Create(opts.OutputPath)
		if err != nil {
			return errs.Config("failed to create output file: " + err.Error())
		}
		defer f.Close()
		out = f
	}

	ctx := context.Background()
	results := batch.Run(ctx, c, ips, actions, batch.RunOptions{
		Concurrency: opts.Concurrency,
		Retries:    opts.Retries,
		DryRun:     opts.DryRun,
	})

	return output.Write(out, results, opts.Format)
}

func resolveActions(action string) []string {
	if action == "all" {
		return endpoint.Actions()
	}
	return []string{action}
}
