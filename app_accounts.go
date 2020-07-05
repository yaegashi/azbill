package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/billing/mgmt/2020-05-01/billing"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppAccounts struct {
	*App
}

func (app *App) AppAccountsCmder() cmder.Cmder {
	return &AppAccounts{App: app}
}

func (app *AppAccounts) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "accounts",
		Aliases:      []string{"a"},
		Short:        "List billing accounts you have access to",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	return cmd
}

func (app *AppAccounts) RunE(cmd *cobra.Command, args []string) error {
	authorizer, err := app.Authorize()
	if err != nil {
		return err
	}

	ctx := context.Background()
	accountsClient := billing.NewAccountsClient("")
	accountsClient.Authorizer = authorizer

	app.Logf("Requesting with %T", accountsClient)

	r, err := accountsClient.ListComplete(ctx, "")
	if err != nil {
		return err
	}

	err = app.Open()
	if err != nil {
		return err
	}
	defer app.Close()

	for r.NotDone() {
		type account billing.Account
		err = app.Marshal(account(r.Value()))
		if err != nil {
			return err
		}
		err = r.NextWithContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
