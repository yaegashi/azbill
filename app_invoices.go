package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2020-05-01-preview/billing"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppInvoices struct {
	*App
	BillingAccount string
	Subscription   string
	StartDate      string
	EndDate        string
}

func (app *App) AppInvoicesCmder() cmder.Cmder {
	return &AppInvoices{App: app}
}

func (app *AppInvoices) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "invoices",
		Aliases:      []string{"i"},
		Short:        "List invoices",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.BillingAccount, "billing-account", "A", "", "billing account")
	cmd.Flags().StringVarP(&app.Subscription, "subscription", "S", "", "subscription")
	cmd.Flags().StringVarP(&app.StartDate, "start", "", "", "start date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&app.EndDate, "end", "", "", "end date (YYYY-MM-DD)")
	return cmd
}

func (app *AppInvoices) RunE(cmd *cobra.Command, args []string) error {
	authorizer, err := app.Authorize()
	if err != nil {
		return err
	}

	ctx := context.Background()
	invoicesClient := billing.NewInvoicesClient(app.Subscription)
	invoicesClient.Authorizer = authorizer

	app.Logf("Requesting with %T", invoicesClient)

	var r billing.InvoiceListResultIterator
	if app.BillingAccount != "" {
		app.Logf("  billing account: %q", app.BillingAccount)
		app.Logf("       start date: %q", app.StartDate)
		app.Logf("         end date: %q", app.EndDate)
		r, err = invoicesClient.ListByBillingAccountComplete(ctx, app.BillingAccount, app.StartDate, app.EndDate)
	} else {
		app.Logf("  subscription: %q", app.Subscription)
		app.Logf("    start date: %q", app.StartDate)
		app.Logf("      end date: %q", app.EndDate)
		r, err = invoicesClient.ListByBillingSubscriptionComplete(ctx, app.StartDate, app.EndDate)
	}
	if err != nil {
		return err
	}

	err = app.Open()
	if err != nil {
		return err
	}
	defer app.Close()

	for r.NotDone() {
		type invoice billing.Invoice
		err = app.Marshal(invoice(r.Value()))
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
