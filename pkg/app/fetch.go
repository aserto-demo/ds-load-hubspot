package app

import (
	"os"

	"github.com/aserto-demo/ds-load-hubspot/pkg/fetch"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
)

type FetchCmd struct {
	ClientID           string `short:"i" help:"Hubspot Client ID" env:"HUBSPOT_CLIENT_ID"`
	ClientSecret       string `short:"s" help:"Hubspot Client Secret" env:"HUBSPOT_CLIENT_SECRET"`
	RefreshToken       string `short:"r" help:"Hubspot Refresh Token" env:"HUBSPOT_REFRESH_TOKEN"`
	PrivateAccessToken string `short:"p" help:"Hubspot Private Access Token" env:"HUBSPOT_PAT"`
	Contacts           bool   `help:"Retrieve Hubspot contacts" env:"HUBSPOT_CONTACTS" default:"false"`
	Companies          bool   `help:"Retrieve Hubspot companies" env:"HUBSPOT_COMPANIES" default:"false"`
}

func (f *FetchCmd) Run(ctx *cc.CommonCtx) error {
	fetcher, err := fetch.New(ctx.Context, f.PrivateAccessToken, f.ClientID, f.ClientSecret, f.RefreshToken)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithOptions(f.Contacts, f.Companies)

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
