package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2020-05-01-preview/billing"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppListAccounts struct {
	*App
}

func (app *App) ListAccountsCmder() cmder.Cmder {
	return &AppListAccounts{App: app}
}

func (app *AppListAccounts) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list-accounts",
		Short:        "List accounts you can have access to",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	return cmd
}

func (app *AppListAccounts) RunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	accountsClient := billing.NewAccountsClient("")
	accountsClient.Authorizer = app.Authorizer

	r, err := accountsClient.ListComplete(ctx, "")
	if err != nil {
		return err
	}

	err = app.Open()
	if err != nil {
		return err
	}
	defer app.Close()

	type account billing.Account
	for r.NotDone() {
		app.JSONMarshal(account(r.Value()))
		err = r.NextWithContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
