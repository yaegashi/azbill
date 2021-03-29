package mapconv

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2020-05-01-preview/billing"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
)

func TestFlatten(t *testing.T) {
	type account billing.Account
	type invoice billing.Invoice
	type subscription subscriptions.Subscription
	type tenant subscriptions.TenantIDDescription
	type legacyUsageDetail consumption.LegacyUsageDetail
	tests := []struct {
		v interface{}
		o bool
		j string
	}{
		{
			v: &P{},
			o: true,
			j: `{"array_int":null,"map_string":null}`,
		},
		{
			v: &P{},
			o: false,
			j: `{"array_int":null,"bool":false,"int":0,"map_string":null,"string":""}`,
		},
		{
			v: &Q{},
			o: true,
			j: `{"A":null,"B":false,"I":0,"M":null,"S":""}`,
		},
		{
			v: &Q{},
			o: false,
			j: `{"A":null,"B":false,"I":0,"M":null,"S":""}`,
		},
		{
			v: &V{},
			o: true,
			j: `{"array_int":null,"bool":false,"int":0,"map_string":null,"p":{"array_int":null,"map_string":null},"q":{"A":null,"B":false,"I":0,"M":null,"S":""},"string":""}`,
		},
		{
			v: &V{},
			o: false,
			j: `{"array_int":null,"bool":false,"int":0,"map_string":null,"p":{"array_int":null,"bool":false,"int":0,"map_string":null,"string":""},"q":{"A":null,"B":false,"I":0,"M":null,"S":""},"string":""}`,
		},
		{
			v: account{},
			o: false,
			j: `{"id":"","name":"","properties":{"accountStatus":"","accountType":"","agreementType":"","billingProfiles":{"hasMoreResults":false,"value":null},"departments":null,"displayName":"","enrollmentAccounts":null,"enrollmentDetails":{"billingCycle":"","channel":"","countryCode":"","currency":"","endDate":"0001-01-01T00:00:00Z","language":"","policies":{"accountOwnerViewCharges":false,"departmentAdminViewCharges":false,"marketplacesEnabled":false,"reservedInstancesEnabled":false},"startDate":"0001-01-01T00:00:00Z","status":""},"hasReadAccess":false,"soldTo":{"addressLine1":"","addressLine2":"","addressLine3":"","city":"","companyName":"","country":"","district":"","email":"","firstName":"","lastName":"","phoneNumber":"","postalCode":"","region":""}},"type":""}`,
		},
		{
			v: invoice{},
			o: false,
			j: `{"id":"","name":"","properties":{"amountDue":{"currency":"","value":0},"azurePrepaymentApplied":{"currency":"","value":0},"billedAmount":{"currency":"","value":0},"billedDocumentId":"","billingProfileDisplayName":"","billingProfileId":"","creditAmount":{"currency":"","value":0},"creditForDocumentId":"","documentType":"","documents":null,"dueDate":"0001-01-01T00:00:00Z","freeAzureCreditApplied":{"currency":"","value":0},"invoiceDate":"0001-01-01T00:00:00Z","invoicePeriodEndDate":"0001-01-01T00:00:00Z","invoicePeriodStartDate":"0001-01-01T00:00:00Z","invoiceType":"","isMonthlyInvoice":false,"payments":null,"purchaseOrderNumber":"","rebillDetails":null,"status":"","subTotal":{"currency":"","value":0},"subscriptionId":"","taxAmount":{"currency":"","value":0},"totalAmount":{"currency":"","value":0}},"type":""}`,
		},
		{
			v: subscription{},
			o: false,
			j: `{"authorizationSource":"","displayName":"","id":"","managedByTenants":null,"state":"","subscriptionId":"","subscriptionPolicies":{"locationPlacementId":"","quotaId":"","spendingLimit":""},"tenantId":""}`,
		},
		{
			v: tenant{},
			o: false,
			j: `{"country":"","countryCode":"","displayName":"","domains":null,"id":"","tenantCategory":"","tenantId":""}`,
		},
		{
			v: legacyUsageDetail{},
			o: false,
			j: `{"id":"","kind":"","name":"","properties":{"accountName":"","accountOwnerId":"","additionalInfo":"","billingAccountId":"","billingAccountName":"","billingCurrency":"","billingPeriodEndDate":"0001-01-01T00:00:00Z","billingPeriodStartDate":"0001-01-01T00:00:00Z","billingProfileId":"","billingProfileName":"","chargeType":"","consumedService":"","cost":"0","costCenter":"","date":"0001-01-01T00:00:00Z","effectivePrice":"0","frequency":"","invoiceSection":"","isAzureCreditEligible":false,"meterDetails":{"meterCategory":"","meterName":"","meterSubCategory":"","serviceFamily":"","unitOfMeasure":""},"meterId":"00000000-0000-0000-0000-000000000000","offerId":"","partNumber":"","planName":"","product":"","productOrderId":"","productOrderName":"","publisherName":"","publisherType":"","quantity":"0","reservationId":"","reservationName":"","resourceGroup":"","resourceId":"","resourceLocation":"","resourceName":"","serviceInfo1":"","serviceInfo2":"","subscriptionId":"","subscriptionName":"","term":"","unitPrice":"0"},"tags":null,"type":""}`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			m, err := Nested(tt.v, tt.o)
			if err != nil {
				t.Fatal(err)
			}
			b, err := json.Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			//t.Log(string(b))
			if tt.j != string(b) {
				t.Errorf("Mismatch\nwant: %s\n got: %s", tt.j, string(b))
			}
		})
	}
}
