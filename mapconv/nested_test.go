package mapconv

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2020-05-01-preview/billing"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
)

func TestNested(t *testing.T) {
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
			j: `{"array_int":null,"bool":false,"int":0,"map_string":null,"p.array_int":null,"p.map_string":null,"q.A":null,"q.B":false,"q.I":0,"q.M":null,"q.S":"","string":""}`,
		},
		{
			v: &V{},
			o: false,
			j: `{"array_int":null,"bool":false,"int":0,"map_string":null,"p.array_int":null,"p.bool":false,"p.int":0,"p.map_string":null,"p.string":"","q.A":null,"q.B":false,"q.I":0,"q.M":null,"q.S":"","string":""}`,
		},
		{
			v: account{},
			o: false,
			j: `{"endDate":"0001-01-01T00:00:00Z","id":"","name":"","properties.accountStatus":"","properties.accountType":"","properties.agreementType":"","properties.billingProfiles.hasMoreResults":false,"properties.billingProfiles.value":null,"properties.departments":null,"properties.displayName":"","properties.enrollmentAccounts":null,"properties.enrollmentDetails.billingCycle":"","properties.enrollmentDetails.channel":"","properties.enrollmentDetails.countryCode":"","properties.enrollmentDetails.currency":"","properties.enrollmentDetails.language":"","properties.enrollmentDetails.policies.accountOwnerViewCharges":false,"properties.enrollmentDetails.policies.departmentAdminViewCharges":false,"properties.enrollmentDetails.policies.marketplacesEnabled":false,"properties.enrollmentDetails.policies.reservedInstancesEnabled":false,"properties.enrollmentDetails.status":"","properties.hasReadAccess":false,"properties.soldTo.addressLine1":"","properties.soldTo.addressLine2":"","properties.soldTo.addressLine3":"","properties.soldTo.city":"","properties.soldTo.companyName":"","properties.soldTo.country":"","properties.soldTo.district":"","properties.soldTo.email":"","properties.soldTo.firstName":"","properties.soldTo.lastName":"","properties.soldTo.phoneNumber":"","properties.soldTo.postalCode":"","properties.soldTo.region":"","startDate":"0001-01-01T00:00:00Z","type":""}`,
		},
		{
			v: invoice{},
			o: false,
			j: `{"dueDate":"0001-01-01T00:00:00Z","id":"","invoiceDate":"0001-01-01T00:00:00Z","invoicePeriodEndDate":"0001-01-01T00:00:00Z","invoicePeriodStartDate":"0001-01-01T00:00:00Z","name":"","properties.amountDue.currency":"","properties.amountDue.value":0,"properties.azurePrepaymentApplied.currency":"","properties.azurePrepaymentApplied.value":0,"properties.billedAmount.currency":"","properties.billedAmount.value":0,"properties.billedDocumentId":"","properties.billingProfileDisplayName":"","properties.billingProfileId":"","properties.creditAmount.currency":"","properties.creditAmount.value":0,"properties.creditForDocumentId":"","properties.documentType":"","properties.documents":null,"properties.freeAzureCreditApplied.currency":"","properties.freeAzureCreditApplied.value":0,"properties.invoiceType":"","properties.isMonthlyInvoice":false,"properties.payments":null,"properties.purchaseOrderNumber":"","properties.rebillDetails":null,"properties.status":"","properties.subTotal.currency":"","properties.subTotal.value":0,"properties.subscriptionId":"","properties.taxAmount.currency":"","properties.taxAmount.value":0,"properties.totalAmount.currency":"","properties.totalAmount.value":0,"type":""}`,
		},
		{
			v: subscription{},
			o: false,
			j: `{"authorizationSource":"","displayName":"","id":"","managedByTenants":null,"state":"","subscriptionId":"","subscriptionPolicies.locationPlacementId":"","subscriptionPolicies.quotaId":"","subscriptionPolicies.spendingLimit":"","tenantId":""}`,
		},
		{
			v: tenant{},
			o: false,
			j: `{"country":"","countryCode":"","displayName":"","domains":null,"id":"","tenantCategory":"","tenantId":""}`,
		},
		{
			v: legacyUsageDetail{},
			o: false,
			j: `{"billingPeriodEndDate":"0001-01-01T00:00:00Z","billingPeriodStartDate":"0001-01-01T00:00:00Z","cost":"0","date":"0001-01-01T00:00:00Z","effectivePrice":"0","id":"","kind":"","meterId":"00000000-0000-0000-0000-000000000000","name":"","properties.accountName":"","properties.accountOwnerId":"","properties.additionalInfo":"","properties.billingAccountId":"","properties.billingAccountName":"","properties.billingCurrency":"","properties.billingProfileId":"","properties.billingProfileName":"","properties.chargeType":"","properties.consumedService":"","properties.costCenter":"","properties.frequency":"","properties.invoiceSection":"","properties.isAzureCreditEligible":false,"properties.meterDetails.meterCategory":"","properties.meterDetails.meterName":"","properties.meterDetails.meterSubCategory":"","properties.meterDetails.serviceFamily":"","properties.meterDetails.unitOfMeasure":"","properties.offerId":"","properties.partNumber":"","properties.planName":"","properties.product":"","properties.productOrderId":"","properties.productOrderName":"","properties.publisherName":"","properties.publisherType":"","properties.reservationId":"","properties.reservationName":"","properties.resourceGroup":"","properties.resourceId":"","properties.resourceLocation":"","properties.resourceName":"","properties.serviceInfo1":"","properties.serviceInfo2":"","properties.subscriptionId":"","properties.subscriptionName":"","properties.term":"","quantity":"0","tags":null,"type":"","unitPrice":"0"}`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			m, err := Flatten(tt.v, tt.o)
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
