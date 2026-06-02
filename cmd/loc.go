package cmd

import (
	"github.com/aiwen/aw-cli/internal/cmdutil"
	ipcmd "github.com/aiwen/aw-cli/cmd/ip"
	"github.com/spf13/cobra"
)

func NewCmdLoc(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "loc <ip>",
		Short: "Query IP geolocation (city, district, street)",
		Long:  "Query geographic location for an IPv4 or IPv6 address. Supports city, district, and street accuracy levels.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "loc"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			opts.Accuracy = ipcmd.ResolveAccuracy(opts)
			return ipcmd.RunIPQuery(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Accuracy, "accuracy", "", "accuracy level: city, district, street (default: city)")
	cmd.Flags().StringVar(&opts.CoordSys, "coordsys", "", "coordinate system (e.g. WGS84, GCJ02)")

	return cmd
}
