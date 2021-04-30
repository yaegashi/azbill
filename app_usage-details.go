package main

import (
	"context"
	"encoding/json"
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
	err = app.Open(ctx)
	if err != nil {
		return err
	}
	defer app.Close(ctx)

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

	var mod func(map[string]interface{}) error
	if app.Flatten {
		mod = func(m1 map[string]interface{}) error {
			if m2, ok := m1["tags"].(map[string]*string); ok {
				s := "{}"
				if m2 != nil {
					if b, err := json.Marshal(m2); err == nil {
						s = string(b)
					}
				}
				m1["tags"] = s
			}
			return nil
		}
	} else {
		mod = func(m1 map[string]interface{}) error {
			if m2, ok := m1["properties"].(map[string]interface{}); ok {
				if m3, ok := m2["additionalInfo"].(string); ok {
					var m4 map[string]interface{}
					if json.Unmarshal([]byte(m3), &m4) == nil {
						m2["additionalInfo"] = m4
					}
				}
			}
			return nil
		}
	}

	for r.NotDone() {
		x := r.Value()
		if v, ok := x.AsLegacyUsageDetail(); ok {
			type LegacyUsageDetail consumption.LegacyUsageDetail
			err = app.Marshal(ctx, (*LegacyUsageDetail)(v), mod)
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
