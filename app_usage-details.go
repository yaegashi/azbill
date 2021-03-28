package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppUsageDetails struct {
	*App
	Scope          string
	BillingAccount string
	BillingPeriod  string
	Subscription   string
	StartDate      string
	EndDate        string
}

func (app *App) AppUsageDetailsCmder() cmder.Cmder {
	return &AppUsageDetails{App: app}
}

func (app *AppUsageDetails) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "usage-details",
		Aliases:      []string{"u"},
		Short:        "List usage details",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.Scope, "scope", "", "", "Scope")
	cmd.Flags().StringVarP(&app.BillingAccount, "billing-account", "A", "", "billing account")
	cmd.Flags().StringVarP(&app.BillingPeriod, "billing-period", "P", "", "billing period")
	cmd.Flags().StringVarP(&app.Subscription, "subscription", "S", "", "subscription")
	cmd.Flags().StringVarP(&app.StartDate, "start", "", "", "start date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&app.EndDate, "end", "", "", "end date (YYYY-MM-DD)")
	return cmd
}

func (app *AppUsageDetails) RunE(cmd *cobra.Command, args []string) error {
	authorizer, err := app.Authorize()
	if err != nil {
		return err
	}

	ctx := context.Background()
	usageDetailsClient := consumption.NewUsageDetailsClient("")
	usageDetailsClient.Authorizer = authorizer

	scope := app.Scope
	if app.BillingAccount != "" {
		scope = filepath.Join(scope, "providers/Microsoft.Billing/billingAccounts", app.BillingAccount)
	}
	if app.Subscription != "" {
		scope = filepath.Join(scope, "subscriptions", app.Subscription)
	}
	if app.BillingPeriod != "" {
		scope = filepath.Join(scope, "providers/Microsoft.Billing/billingPeriods", app.BillingPeriod)
	}
	if scope == "" {
		return fmt.Errorf("no scope specified")
	}

	expand := "properties/additionalInfo,properties/meterDetails"
	filter := ""
	if app.StartDate != "" && app.EndDate != "" {
		filter = fmt.Sprintf("properties/usageStart eq '%s' and properties/usageEnd eq '%s'", app.StartDate, app.EndDate)
	}

	app.Logf("Requesting with %T", usageDetailsClient)
	app.Logf("   scope: %q", scope)
	app.Logf("  filter: %q", filter)

	r, err := usageDetailsClient.ListComplete(ctx, scope, expand, filter, "", nil, "")
	if err != nil {
		return err
	}

	err = app.Open()
	if err != nil {
		return err
	}
	defer app.Close()

	for r.NotDone() {
		x := r.Value()
		if v, ok := x.AsLegacyUsageDetail(); ok {
			type LegacyUsageDetail consumption.LegacyUsageDetail
			err = app.Marshal((*LegacyUsageDetail)(v))
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unexpected type %T", x)
		}
		err = r.NextWithContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
