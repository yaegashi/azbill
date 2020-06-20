package main

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yaegashi/azbill/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	cmder "github.com/yaegashi/cobra-cmder"
)

type AppListUsageDetails struct {
	*App
	Scope          string
	BillingAccount string
	BillingPeriod  string
	Subscription   string
	StartDate      string
	EndDate        string
}

func (app *App) ListUsageDetailsCmder() cmder.Cmder {
	return &AppListUsageDetails{App: app}
}

func (app *AppListUsageDetails) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list-usage-details",
		Short:        "List usage details",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.Scope, "scope", "", "", "Scope")
	cmd.Flags().StringVarP(&app.BillingAccount, "billing-account", "A", "", "Billing account")
	cmd.Flags().StringVarP(&app.BillingPeriod, "billing-period", "P", "", "Billing period")
	cmd.Flags().StringVarP(&app.Subscription, "subscription", "S", "", "Subscription")
	cmd.Flags().StringVarP(&app.StartDate, "start-date", "", "", "Start date")
	cmd.Flags().StringVarP(&app.EndDate, "end-date", "", "", "Start date")
	return cmd
}

func (app *AppListUsageDetails) RunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	usageDetailsClient := consumption.NewUsageDetailsClient("")
	usageDetailsClient.Authorizer = app.Authorizer

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
		return fmt.Errorf("No scope specified")
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

	var keys []string
	if app.Format == "csv" {
		m, _ := flattenToMap(consumption.LegacyUsageDetail{}, false)
		for key := range m {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		app.CSVMarshal(keys)
	}

	for r.NotDone() {
		x, _ := r.Value().AsLegacyUsageDetail()
		if app.Format == "csv" {
			row := []string{}
			m, err := flattenToMap(x, true)
			if err == nil {
				for _, key := range keys {
					col := ""
					if val, ok := m[key]; ok {
						if tags, ok := val.(map[string]*string); ok {
							if tags != nil {
								b, _ := json.Marshal(val)
								col = string(b)
							}
						} else {
							col = fmt.Sprint(val)
						}
					}
					row = append(row, col)
				}
				err = app.CSVMarshal(row)
			}
		} else if app.Flatten {
			m, err := flattenToMap(x, true)
			if err == nil {
				err = app.JSONMarshal(m)
			}
		} else {
			type LegacyUsageDetail consumption.LegacyUsageDetail
			err = app.JSONMarshal((*LegacyUsageDetail)(x))
		}
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
