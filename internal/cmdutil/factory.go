package cmdutil

import (
	"net/http"

	"github.com/aiwen/aw-cli/internal/client"
	"github.com/aiwen/aw-cli/internal/core"
)

type GlobalOptions struct {
	ConfigPath string
	BaseURL    string
	APIKey     string
	Timeout    string
	Format     string
	JQ         string
	DryRun     bool
	Verbose    bool
}

type Factory struct {
	IO      IOStreams
	Options *GlobalOptions
	Config  core.CliConfig
}

func NewFactory(io IOStreams, opts *GlobalOptions) *Factory {
	return &Factory{IO: io, Options: opts}
}

func (f *Factory) ResolveConfig() (core.CliConfig, error) {
	cfg, err := core.Resolve(core.Overrides{
		ConfigPath: f.Options.ConfigPath,
		BaseURL:    f.Options.BaseURL,
		APIKey:     f.Options.APIKey,
		Timeout:    f.Options.Timeout,
	})
	if err != nil {
		return core.CliConfig{}, err
	}
	f.Config = cfg
	return cfg, nil
}

func (f *Factory) Client() (*client.AiwenClient, error) {
	cfg, err := f.ResolveConfig()
	if err != nil {
		return nil, err
	}
	return &client.AiwenClient{
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Channel: cfg.Channel,
		HTTP: &http.Client{
			Timeout: core.ParseTimeout(cfg.Timeout),
		},
		ErrOut: f.IO.ErrOut,
	}, nil
}
