package consumptionapi

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/yaegashi/azbill/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
)

// UsageDetailsClientAPI contains the set of methods on the UsageDetailsClient type.
type UsageDetailsClientAPI interface {
	List(ctx context.Context, scope string, expand string, filter string, skiptoken string, top *int32, metric consumption.Metrictype) (result consumption.UsageDetailsListResultPage, err error)
	ListComplete(ctx context.Context, scope string, expand string, filter string, skiptoken string, top *int32, metric consumption.Metrictype) (result consumption.UsageDetailsListResultIterator, err error)
}

var _ UsageDetailsClientAPI = (*consumption.UsageDetailsClient)(nil)

// MarketplacesClientAPI contains the set of methods on the MarketplacesClient type.
type MarketplacesClientAPI interface {
	List(ctx context.Context, scope string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListComplete(ctx context.Context, scope string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
}

var _ MarketplacesClientAPI = (*consumption.MarketplacesClient)(nil)

// BudgetsClientAPI contains the set of methods on the BudgetsClient type.
type BudgetsClientAPI interface {
	CreateOrUpdate(ctx context.Context, scope string, budgetName string, parameters consumption.Budget) (result consumption.Budget, err error)
	Delete(ctx context.Context, scope string, budgetName string) (result autorest.Response, err error)
	Get(ctx context.Context, scope string, budgetName string) (result consumption.Budget, err error)
	List(ctx context.Context, scope string) (result consumption.BudgetsListResultPage, err error)
	ListComplete(ctx context.Context, scope string) (result consumption.BudgetsListResultIterator, err error)
}

var _ BudgetsClientAPI = (*consumption.BudgetsClient)(nil)

// TagsClientAPI contains the set of methods on the TagsClient type.
type TagsClientAPI interface {
	Get(ctx context.Context, scope string) (result consumption.TagsResult, err error)
}

var _ TagsClientAPI = (*consumption.TagsClient)(nil)

// ChargesClientAPI contains the set of methods on the ChargesClient type.
type ChargesClientAPI interface {
	List(ctx context.Context, scope string, startDate string, endDate string, filter string, apply string) (result consumption.ChargesListResult, err error)
}

var _ ChargesClientAPI = (*consumption.ChargesClient)(nil)

// BalancesClientAPI contains the set of methods on the BalancesClient type.
type BalancesClientAPI interface {
	GetByBillingAccount(ctx context.Context, billingAccountID string) (result consumption.Balance, err error)
	GetForBillingPeriodByBillingAccount(ctx context.Context, billingAccountID string, billingPeriodName string) (result consumption.Balance, err error)
}

var _ BalancesClientAPI = (*consumption.BalancesClient)(nil)

// ReservationsSummariesClientAPI contains the set of methods on the ReservationsSummariesClient type.
type ReservationsSummariesClientAPI interface {
	List(ctx context.Context, scope string, grain consumption.Datagrain, startDate string, endDate string, filter string) (result consumption.ReservationSummariesListResultPage, err error)
	ListComplete(ctx context.Context, scope string, grain consumption.Datagrain, startDate string, endDate string, filter string) (result consumption.ReservationSummariesListResultIterator, err error)
	ListByReservationOrder(ctx context.Context, reservationOrderID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultPage, err error)
	ListByReservationOrderComplete(ctx context.Context, reservationOrderID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultIterator, err error)
	ListByReservationOrderAndReservation(ctx context.Context, reservationOrderID string, reservationID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultPage, err error)
	ListByReservationOrderAndReservationComplete(ctx context.Context, reservationOrderID string, reservationID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultIterator, err error)
}

var _ ReservationsSummariesClientAPI = (*consumption.ReservationsSummariesClient)(nil)

// ReservationsDetailsClientAPI contains the set of methods on the ReservationsDetailsClient type.
type ReservationsDetailsClientAPI interface {
	List(ctx context.Context, scope string, startDate string, endDate string, filter string) (result consumption.ReservationDetailsListResultPage, err error)
	ListComplete(ctx context.Context, scope string, startDate string, endDate string, filter string) (result consumption.ReservationDetailsListResultIterator, err error)
	ListByReservationOrder(ctx context.Context, reservationOrderID string, filter string) (result consumption.ReservationDetailsListResultPage, err error)
	ListByReservationOrderComplete(ctx context.Context, reservationOrderID string, filter string) (result consumption.ReservationDetailsListResultIterator, err error)
	ListByReservationOrderAndReservation(ctx context.Context, reservationOrderID string, reservationID string, filter string) (result consumption.ReservationDetailsListResultPage, err error)
	ListByReservationOrderAndReservationComplete(ctx context.Context, reservationOrderID string, reservationID string, filter string) (result consumption.ReservationDetailsListResultIterator, err error)
}

var _ ReservationsDetailsClientAPI = (*consumption.ReservationsDetailsClient)(nil)

// ReservationRecommendationsClientAPI contains the set of methods on the ReservationRecommendationsClient type.
type ReservationRecommendationsClientAPI interface {
	List(ctx context.Context, scope string, filter string) (result consumption.ReservationRecommendationsListResultPage, err error)
	ListComplete(ctx context.Context, scope string, filter string) (result consumption.ReservationRecommendationsListResultIterator, err error)
}

var _ ReservationRecommendationsClientAPI = (*consumption.ReservationRecommendationsClient)(nil)

// ReservationRecommendationDetailsClientAPI contains the set of methods on the ReservationRecommendationDetailsClient type.
type ReservationRecommendationDetailsClientAPI interface {
	Get(ctx context.Context, scope string) (result consumption.ReservationRecommendationDetailsModel, err error)
}

var _ ReservationRecommendationDetailsClientAPI = (*consumption.ReservationRecommendationDetailsClient)(nil)

// ReservationTransactionsClientAPI contains the set of methods on the ReservationTransactionsClient type.
type ReservationTransactionsClientAPI interface {
	List(ctx context.Context, billingAccountID string, filter string) (result consumption.ReservationTransactionsListResultPage, err error)
	ListComplete(ctx context.Context, billingAccountID string, filter string) (result consumption.ReservationTransactionsListResultIterator, err error)
	ListByBillingProfile(ctx context.Context, billingAccountID string, billingProfileID string, filter string) (result consumption.ModernReservationTransactionsListResultPage, err error)
	ListByBillingProfileComplete(ctx context.Context, billingAccountID string, billingProfileID string, filter string) (result consumption.ModernReservationTransactionsListResultIterator, err error)
}

var _ ReservationTransactionsClientAPI = (*consumption.ReservationTransactionsClient)(nil)

// PriceSheetClientAPI contains the set of methods on the PriceSheetClient type.
type PriceSheetClientAPI interface {
	Get(ctx context.Context, expand string, skiptoken string, top *int32) (result consumption.PriceSheetResult, err error)
	GetByBillingPeriod(ctx context.Context, billingPeriodName string, expand string, skiptoken string, top *int32) (result consumption.PriceSheetResult, err error)
}

var _ PriceSheetClientAPI = (*consumption.PriceSheetClient)(nil)

// ForecastsClientAPI contains the set of methods on the ForecastsClient type.
type ForecastsClientAPI interface {
	List(ctx context.Context, filter string) (result consumption.ForecastsListResult, err error)
}

var _ ForecastsClientAPI = (*consumption.ForecastsClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result consumption.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result consumption.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*consumption.OperationsClient)(nil)

// AggregatedCostClientAPI contains the set of methods on the AggregatedCostClient type.
type AggregatedCostClientAPI interface {
	GetByManagementGroup(ctx context.Context, managementGroupID string, filter string) (result consumption.ManagementGroupAggregatedCostResult, err error)
	GetForBillingPeriodByManagementGroup(ctx context.Context, managementGroupID string, billingPeriodName string) (result consumption.ManagementGroupAggregatedCostResult, err error)
}

var _ AggregatedCostClientAPI = (*consumption.AggregatedCostClient)(nil)

// EventsClientAPI contains the set of methods on the EventsClient type.
type EventsClientAPI interface {
	List(ctx context.Context, billingAccountID string, billingProfileID string, startDate string, endDate string) (result consumption.EventsPage, err error)
	ListComplete(ctx context.Context, billingAccountID string, billingProfileID string, startDate string, endDate string) (result consumption.EventsIterator, err error)
}

var _ EventsClientAPI = (*consumption.EventsClient)(nil)

// LotsClientAPI contains the set of methods on the LotsClient type.
type LotsClientAPI interface {
	List(ctx context.Context, billingAccountID string, billingProfileID string) (result consumption.LotsPage, err error)
	ListComplete(ctx context.Context, billingAccountID string, billingProfileID string) (result consumption.LotsIterator, err error)
}

var _ LotsClientAPI = (*consumption.LotsClient)(nil)

// CreditsClientAPI contains the set of methods on the CreditsClient type.
type CreditsClientAPI interface {
	Get(ctx context.Context, billingAccountID string, billingProfileID string) (result consumption.CreditSummary, err error)
}

var _ CreditsClientAPI = (*consumption.CreditsClient)(nil)
