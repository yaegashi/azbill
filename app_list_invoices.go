package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2020-05-01-preview/billing"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppListInvoices struct {
	*App
	StartDate      string
	EndDate        string
	BillingAccount string
}

func (app *App) ListInvoicesCmder() cmder.Cmder {
	return &AppListInvoices{App: app}
}

func (app *AppListInvoices) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list-invoices",
		Short:        "List invoices",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.StartDate, "start-date", "S", "", "Start date")
	cmd.Flags().StringVarP(&app.EndDate, "end-date", "E", "", "End date")
	cmd.Flags().StringVarP(&app.BillingAccount, "billing-account", "A", "", "Billing account")
	return cmd
}

func (app *AppListInvoices) RunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	invoicesClient := billing.NewInvoicesClient("")
	invoicesClient.Authorizer = app.Authorizer
	r, err := invoicesClient.ListByBillingAccountComplete(ctx, app.BillingAccount, app.StartDate, app.EndDate)
	if err != nil {
		return err
	}

	err = app.Open()
	if err != nil {
		return err
	}
	defer app.Close()

	type invoice billing.Invoice
	for r.NotDone() {
		app.JSONMarshal(invoice(r.Value()))
		err = r.NextWithContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
