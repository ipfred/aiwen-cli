package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aiwen/aw-cli/errs"
	"github.com/aiwen/aw-cli/internal/build"
	"github.com/aiwen/aw-cli/internal/cmdutil"
	"github.com/aiwen/aw-cli/internal/output"
	"github.com/spf13/cobra"
)

func Execute() int {
	f := cmdutil.NewFactory(cmdutil.SystemIOStreams(), &cmdutil.GlobalOptions{})
	rootCmd := NewRootCmd(f)
	ctx := context.Background()
	rootCmd.SetContext(ctx)
	if err := rootCmd.Execute(); err != nil {
		return handleRootError(f, err)
	}
	return 0
}

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	root := &cobra.Command{
		Use:           "aw-cli",
		Short:         "AIWEN/IPPlus360 IP intelligence query CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       build.Version,
	}

	SetGlobalFlags(root, f.Options)

	root.AddCommand(
		NewCmdLoc(f),
		NewCmdCurrent(f),
		NewCmdScene(f),
		NewCmdWhois(f),
		NewCmdASN(f),
		NewCmdHost(f),
		NewCmdRisk(f),
		NewCmdIdentity(f),
		NewCmdIndustry(f),
		NewCmdBatch(f),
		NewCmdConfig(f),
		NewCmdCompletion(root),
	)

	return root
}

func SetGlobalFlags(cmd *cobra.Command, opts *cmdutil.GlobalOptions) {
	cmd.PersistentFlags().StringVar(&opts.ConfigPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&opts.BaseURL, "base-url", "", "API base URL (default: https://api.ipplus360.com)")
	cmd.PersistentFlags().StringVar(&opts.APIKey, "api-key", "", "API key (overrides config and env)")
	cmd.PersistentFlags().StringVar(&opts.Timeout, "timeout", "", "HTTP timeout duration (default: 10s)")
	cmd.PersistentFlags().StringVar(&opts.Format, "format", output.FormatJSON, "output format: json, ndjson, table, csv")
	cmd.PersistentFlags().StringVar(&opts.JQ, "jq", "", "jq-style filter expression (e.g. .data.country)")
	cmd.PersistentFlags().BoolVar(&opts.DryRun, "dry-run", false, "print request without calling upstream")
	cmd.PersistentFlags().BoolVar(&opts.Verbose, "verbose", false, "verbose output; secrets are redacted")
}

func handleRootError(f *cmdutil.Factory, err error) int {
	if exitErr, ok := err.(*errs.ExitError); ok {
		writeError(f.IO.ErrOut, exitErr.Problem, f.Options.Format)
		return exitErr.Code
	}
	problem := errs.Problem{
		Type:    errs.TypeInternal,
		Message: err.Error(),
	}
	writeError(f.IO.ErrOut, problem, f.Options.Format)
	return 1
}

func writeError(w io.Writer, p errs.Problem, format string) {
	if format == output.FormatTable || format == output.FormatCSV {
		fmt.Fprintln(w, p.Message)
		return
	}
	env := output.Envelope{
		OK:    false,
		Error: p,
	}
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		fmt.Fprintln(w, p.Message)
		return
	}
	fmt.Fprintln(w, string(data))
}
