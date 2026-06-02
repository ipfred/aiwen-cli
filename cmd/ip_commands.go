package cmd

import (
	"github.com/aiwen/aw-cli/internal/cmdutil"
	ipcmd "github.com/aiwen/aw-cli/cmd/ip"
	"github.com/spf13/cobra"
)

func NewCmdWhois(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "whois <ip>",
		Short: "Query IP WHOIS registration information",
		Long:  "Query WHOIS registration data for an IPv4 address, including registered network, organization, and contacts.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "whois"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "whois"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	return cmd
}

func NewCmdASN(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "asn <ip>",
		Short: "Query AS number and AS WHOIS for an IP",
		Long:  "Map an IPv4 address to its autonomous system number and AS WHOIS information.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "asn"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "asn"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	return cmd
}

func NewCmdHost(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "host <ip>",
		Short: "Query IP host information (AS name, ISP, organization)",
		Long:  "Query host information for an IPv4 address, including AS number, AS name, ISP, and owning organization.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "host"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "host"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	return cmd
}

func NewCmdRisk(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "risk <ip>",
		Short: "Query IP risk portrait (VPN, proxy, Tor, datacenter)",
		Long:  "Query risk portrait for an IPv4 address, including VPN, proxy, Tor, datacenter detection, and risk scoring.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "risk"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "risk"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	return cmd
}

func NewCmdIdentity(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "identity <ip>",
		Short: "Query IP identity (real human vs bot traffic)",
		Long:  "Analyze an IPv4 address to determine the probability of real human traffic vs machine/bot traffic and second-use probability.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "identity"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "identity"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	return cmd
}

func NewCmdIndustry(f *cmdutil.Factory) *cobra.Command {
	opts := &ipcmd.IPQueryOptions{Factory: f}

	cmd := &cobra.Command{
		Use:   "industry <ip>",
		Short: "Query IP industry classification",
		Long:  "Classify an IPv4 address by its associated industry sector.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ctx = cmd.Context()
			opts.Action = "industry"
			opts.IP = args[0]
			opts.Format = f.Options.Format
			opts.JqExpr = f.Options.JQ
			opts.DryRun = f.Options.DryRun
			if err := ipcmd.ValidateIPForAction(opts.IP, "industry"); err != nil {
				return err
			}
			return ipcmd.RunIPQuery(opts)
		},
	}

	return cmd
}
