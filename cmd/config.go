package cmd

import (
	"fmt"

	"github.com/aiwen/aw-cli/internal/cmdutil"
	"github.com/aiwen/aw-cli/internal/core"
	"github.com/aiwen/aw-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  "Manage the aw-cli CLI configuration file. Use subcommands to initialize, show, or set configuration values.",
	}

	cmd.AddCommand(NewCmdConfigInit(f))
	cmd.AddCommand(NewCmdConfigShow(f))
	cmd.AddCommand(NewCmdConfigSet(f))

	return cmd
}

func NewCmdConfigInit(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a default configuration file",
		Long:  "Create a default configuration file at ~/.aw-cli/config.json with sensible defaults.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := core.DefaultConfigPath()
			if err != nil {
				return err
			}
			cfg := core.CliConfig{
				BaseURL:      core.DefaultBaseURL,
				Channel:      core.DefaultChannel,
				Timeout:      core.DefaultTimeout.String(),
				IPv4Accuracy: core.DefaultAccuracy,
				IPv6Accuracy: core.DefaultAccuracy,
			}
			if err := core.Write(path, cfg); err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "Configuration written to %s\n", path)
			return nil
		},
	}

	return cmd
}

func NewCmdConfigShow(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display current configuration",
		Long:  "Show the resolved configuration values, including defaults, config file, and environment variable overrides.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.ResolveConfig()
			if err != nil {
				return err
			}
			redacted := cfg
			redacted.APIKey = core.RedactSecret(cfg.APIKey)
			return output.Write(f.IO.Out, redacted, f.Options.Format)
		},
	}

	return cmd
}

func NewCmdConfigSet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  "Set a configuration value in the config file. Supported keys: base_url, api_key, channel, timeout, ipv4_accuracy, ipv6_accuracy.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			cfg, err := core.Load(f.Options.ConfigPath)
			if err != nil {
				return err
			}
			switch key {
			case "base_url":
				cfg.BaseURL = value
			case "api_key":
				cfg.APIKey = value
			case "channel":
				cfg.Channel = value
			case "timeout":
				cfg.Timeout = value
			case "ipv4_accuracy":
				cfg.IPv4Accuracy = value
			case "ipv6_accuracy":
				cfg.IPv6Accuracy = value
			default:
				return fmt.Errorf("unknown config key: %s; valid keys: base_url, api_key, channel, timeout, ipv4_accuracy, ipv6_accuracy", key)
			}
			if err := core.Write(f.Options.ConfigPath, cfg); err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "Set %s\n", key)
			return nil
		},
	}

	return cmd
}
