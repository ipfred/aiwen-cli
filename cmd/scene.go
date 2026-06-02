package cmd

import (
	"github.com/aiwen/aw-cli/internal/cmdutil"
	ipcmd "github.com/aiwen/aw-cli/cmd/ip"
	"github.com/spf13/cobra"
)

func NewCmdScene(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "scene <ip>",
		Short: "Query IP usage scene (residential, datacenter, CDN, etc.)",
		Long:  "Query the usage scene of an IP address, such as residential, datacenter, CDN, Anycast, satellite, etc.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "scene"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "scene"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Lang, "lang", "", "response language (default: cn)")

	return cmd
}
