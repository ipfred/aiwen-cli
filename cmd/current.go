package cmd

import (
	"github.com/aiwen/aw-cli/internal/cmdutil"
	ipcmd "github.com/aiwen/aw-cli/cmd/ip"
	"github.com/spf13/cobra"
)

func NewCmdCurrent(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Query geolocation for the current network IP",
		Long:  "Detect your current public IP address and query its geolocation information.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "current"
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			return ipcmd.RunCurrentQuery(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Accuracy, "accuracy", "", "accuracy level: city, district, street (default: city)")
	cmd.Flags().StringVar(&opts.CoordSys, "coordsys", "", "coordinate system (e.g. WGS84, GCJ02)")

	return cmd
}
