package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppTenants struct {
	*App
	StartDate      string
	EndDate        string
	BillingAccount string
}

func (app *App) AppTenantsCmder() cmder.Cmder {
	return &AppTenants{App: app}
}

func (app *AppTenants) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "tenants",
		Aliases:      []string{"t"},
		Short:        "List tenants",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	return cmd
}

func (app *AppTenants) RunE(cmd *cobra.Command, args []string) error {
	authorizer, err := app.Authorize()
	if err != nil {
		return err
	}

	ctx := context.Background()
	tenantsClient := subscriptions.NewTenantsClient()
	tenantsClient.Authorizer = authorizer

	app.Logf("Requesting with %T", tenantsClient)

	r, err := tenantsClient.ListComplete(ctx)
	if err != nil {
		return err
	}

	err = app.Open()
	if err != nil {
		return err
	}
	defer app.Close()

	for r.NotDone() {
		type tenant subscriptions.TenantIDDescription
		err = app.Marshal(tenant(r.Value()))
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
