package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppSubscriptions struct {
	*App
	StartDate      string
	EndDate        string
	BillingAccount string
}

func (app *App) AppSubscriptionsCmder() cmder.Cmder {
	return &AppSubscriptions{App: app}
}

func (app *AppSubscriptions) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "subscriptions",
		Aliases:      []string{"s"},
		Short:        "List subscriptions",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	return cmd
}

func (app *AppSubscriptions) RunE(cmd *cobra.Command, args []string) error {
	authorizer, err := app.Authorize()
	if err != nil {
		return err
	}

	ctx := context.Background()
	subscriptionsClient := subscriptions.NewClient()
	subscriptionsClient.Authorizer = authorizer

	app.Logf("Requesting with %T", subscriptionsClient)

	r, err := subscriptionsClient.ListComplete(ctx)
	if err != nil {
		return err
	}

	err = app.Open(ctx)
	if err != nil {
		return err
	}
	defer app.Close(ctx)

	for r.NotDone() {
		type subscription subscriptions.Subscription
		err = app.Marshal(ctx, subscription(r.Value()))
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
