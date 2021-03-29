# azbill

## Introduction

azbill is a cross-platform CLI tool to bulk export Azure billing data in CSV or JSONL format
which is useful for the offline cost analysis.

Currently azbill supports the following [billing account types](https://docs.microsoft.com/en-us/azure/cost-management-billing/manage/view-all-accounts):

- Microsoft Online Services Program
- Enterprise Agreement

## Usage

```console
$ azbill -h
Azure billing data exporter

Usage:
  azbill [command]

Available Commands:
  accounts      List billing accounts you have access to
  help          Help about any command
  invoices      List invoices
  login         Force auth-dev login
  subscriptions List subscriptions
  tenants       List tenants
  usage-details List usage details

Flags:
      --auth string               auth source [dev,env,file,cli] (env:AZBILL_AUTH, default:dev)
      --auth-dev string           auth dev store (env:AZBILL_AUTH_DEV, default:auth_dev.json)
      --auth-file string          auth file store (env:AZBILL_AUTH_FILE, default:auth_file.json)
      --client string             Azure client (env:AZURE_CLIENT_ID, default:4a034c56-da44-48ce-90db-039a408974bd)
      --config-dir string         config dir (env:AZBILL_CONFIG_DIR, default:~/.azbill)
      --format string             output format [csv,json,flatten,pretty] (env:AZBILL_FORMAT, default:csv)
  -h, --help                      help for azbill
      --mongo-collection string   output MongoDB collection
      --mongo-db string           output MongoDB database
      --mongo-drop                drop the existing MongoDB collection
      --mongo-uri string          output MongoDB URI
  -o, --output string             output file path
  -q, --quiet                     quiet
      --tenant string             Azure tenant (env:AZURE_TENANT_ID, default:common)
  -v, --version                   version for azbill

Use "azbill [command] --help" for more information about a command.
```

## Authentication

You can sign in using [device authorization grant flow](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code) by `azbill login`:

```console
$ azbill login
2020/07/06 23:05:05 To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code HL94ZL7Y8 to authenticate.
2020/07/06 23:05:21 Saving auth-dev token in /Users/yaegashi/.azbill/auth_dev.json
```

To sign in with other tenant than your home tenant:

```console
$ azbill login --tenant l0wdev.onmicrosoft.com
```

The persistent auth token is saved in `~/.azbill/auth_dev.json` per default.
You can specify a custom location to save it with `--auth-dev`.
It's a relative path in the config dir specifyed by `--config-dir` (default: `~/.azbill`).
You can also pass an Azure Blob Storage URL with SAS, which is especially useful in the CI/CD environment.

[Other authentication methods supported by Azure SDK for Go](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization) are also available.  You can select the preferred method by `--auth`.
If you've already signed in with the Azure CLI, `--auth cli` would be most useful.
You can use an auth file generated by the Azure CLI by `--auth file` and `--auth-file` to specify its location.

## Output formats

### json

With `--format json`, the output is a series of JSON objects like the following:

```json
{
  "properties": {
    "billingPeriodStartDate": "2020-05-03T00:00:00Z",
    "billingPeriodEndDate": "2020-06-02T00:00:00Z",
    "billingProfileId": "/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
    "billingProfileName": "Pay-As-You-Go",
    "subscriptionId": "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
    "subscriptionName": "Pay-As-You-Go",
    "date": "2020-06-01T00:00:00Z",
    "meterId": "d54686f0-77ff-43f3-9e7c-2099030d32a7",
    "meterDetails": {
      "meterName": "Public Queries",
      "meterCategory": "Azure DNS",
      "meterSubCategory": "",
      "unitOfMeasure": "10000000"
    },
    "quantity": "0.000997",
    "effectivePrice": "55.3250345781466",
    "cost": "0.055159059474412",
    "unitPrice": "0",
    "billingCurrency": "JPY",
    "resourceLocation": "Unknown",
    "consumedService": "Microsoft.Network",
    "resourceId": "/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/resourceGroups/dns/providers/Microsoft.Network/dnszones/l0w.dev",
    "resourceName": "l0w.dev",
    "resourceGroup": "dns",
    "offerId": "MS-AZR-0003P",
    "isAzureCreditEligible": false,
    "publisherType": "Azure",
    "chargeType": "Usage",
    "frequency": "UsageBased"
  },
  "id": "/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/providers/Microsoft.Billing/billingPeriods/202006/providers/Microsoft.Consumption/usageDetails/de545724-82e8-f099-bc3b-53c5f0ae1711",
  "name": "de545724-82e8-f099-bc3b-53c5f0ae1711",
  "type": "Microsoft.Consumption/usageDetails",
  "tags": null,
  "kind": "legacy"
}
```

Note that numerical values are expressed in string to avoid floating point rounding errors.

The actual output is in [JSONL format](http://jsonlines.org) with `--format json`,
where the JSON objects contain no whitespaces and they are delimited by newlines.
You can get pretty-printed JSON objects like above with `--format json,pretty`.

### flatten

With `--format flatten`, the output is a series of JSON objects like the following:

```json
{
  "id": "/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/providers/Microsoft.Billing/billingPeriods/202006/providers/Microsoft.Consumption/usageDetails/de545724-82e8-f099-bc3b-53c5f0ae1711",
  "kind": "legacy",
  "name": "de545724-82e8-f099-bc3b-53c5f0ae1711",
  "properties.billingCurrency": "JPY",
  "properties.billingPeriodEndDate": "2020-06-02T00:00:00Z",
  "properties.billingPeriodStartDate": "2020-05-03T00:00:00Z",
  "properties.billingProfileId": "/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
  "properties.billingProfileName": "Pay-As-You-Go",
  "properties.chargeType": "Usage",
  "properties.consumedService": "Microsoft.Network",
  "properties.cost": "0.055159059474412",
  "properties.date": "2020-06-01T00:00:00Z",
  "properties.effectivePrice": "55.3250345781466",
  "properties.frequency": "UsageBased",
  "properties.isAzureCreditEligible": false,
  "properties.meterDetails.meterCategory": "Azure DNS",
  "properties.meterDetails.meterName": "Public Queries",
  "properties.meterDetails.meterSubCategory": "",
  "properties.meterDetails.unitOfMeasure": "10000000",
  "properties.meterId": "d54686f0-77ff-43f3-9e7c-2099030d32a7",
  "properties.offerId": "MS-AZR-0003P",
  "properties.publisherType": "Azure",
  "properties.quantity": "0.000997",
  "properties.resourceGroup": "dns",
  "properties.resourceId": "/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/resourceGroups/dns/providers/Microsoft.Network/dnszones/l0w.dev",
  "properties.resourceLocation": "Unknown",
  "properties.resourceName": "l0w.dev",
  "properties.subscriptionId": "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
  "properties.subscriptionName": "Pay-As-You-Go",
  "properties.unitPrice": "0",
  "tags": null,
  "type": "Microsoft.Consumption/usageDetails"
}
```

It's like `json` but all nested objects are flatten to form a single object.

### csv

With `--format csv`, the output is CSV records like the following:

```csv
id,kind,name,properties.accountName,properties.accountOwnerId,properties.additionalInfo,properties.billingAccountId,properties.billingAccountName,properties.billingCurrency,properties.billingPeriodEndDate,properties.billingPeriodStartDate,properties.billingProfileId,properties.billingProfileName,properties.chargeType,properties.consumedService,properties.cost,properties.costCenter,properties.date,properties.effectivePrice,properties.frequency,properties.invoiceSection,properties.isAzureCreditEligible,properties.meterDetails.meterCategory,properties.meterDetails.meterName,properties.meterDetails.meterSubCategory,properties.meterDetails.serviceFamily,properties.meterDetails.unitOfMeasure,properties.meterId,properties.offerId,properties.partNumber,properties.planName,properties.product,properties.productOrderId,properties.productOrderName,properties.publisherName,properties.publisherType,properties.quantity,properties.reservationId,properties.reservationName,properties.resourceGroup,properties.resourceId,properties.resourceLocation,properties.resourceName,properties.serviceInfo1,properties.serviceInfo2,properties.subscriptionId,properties.subscriptionName,properties.term,properties.unitPrice,tags,type
/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/providers/Microsoft.Billing/billingPeriods/202006/providers/Microsoft.Consumption/usageDetails/de545724-82e8-f099-bc3b-53c5f0ae1711,legacy,de545724-82e8-f099-bc3b-53c5f0ae1711,,,,,,JPY,2020-06-02T00:00:00Z,2020-05-03T00:00:00Z,/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX,Pay-As-You-Go,Usage,Microsoft.Network,0.055159059474412,,2020-06-01T00:00:00Z,55.3250345781466,UsageBased,,false,Azure DNS,Public Queries,,,10000000,d54686f0-77ff-43f3-9e7c-2099030d32a7,MS-AZR-0003P,,,,,,,Azure,0.000997,,,dns,/subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/resourceGroups/dns/providers/Microsoft.Network/dnszones/l0w.dev,Unknown,l0w.dev,,,XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX,Pay-As-You-Go,,0,{},Microsoft.Consumption/usageDetails
```

The first line is the CSV header based on JSON object keys of `flatten` format.
The attributes omit due to empty values in `flatten` format are all present in it.
The encoding is UTF-8 with BOM, the line ending is CRLF.
You should be able to directly open it with Microsoft Excel.

## Output destination

### File or standard output

azbill makes the output into a file specified by `--output`.
If `--output` is not specified, it writes to the standard output.

`--format` defaults to `csv`.

### MongoDB

With `--mongo-uri`, azbill connects to the MongoDB server and makes the output there.
`--output` is ignored.

It also requires `--mongo-db` and `--mongo-collection` for the MongoDB output to work.
Specify `--mongo-drop` to drop the existing collection before writing.

`--format` should be `json` (default) or `flatten`.
It writes a document for every JSON object.  

## Examples

List [billing accounts](https://docs.microsoft.com/en-us/azure/cost-management-billing/manage/view-all-accounts) you have access to in CSV format:

```console
$ azbill accounts
```

List tenants in JSONL format:

```console
$ azbill tenants --format json
```

List subscriptions in flatten JSONL format:

```console
$ azbill subscriptions --format flatten
```

List invoices of the first half of 2020 for subscription XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX in pretty JSON format:

```console
$ azbill invoices --format pretty -S XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX --start 2020-06-01 --end 2020-06-30
```

For Microsoft Online Service Program accounts: List usage details of June 2020 for subscription XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX into usage.jsonl in flatten and pretty JSONL format:

```console
$ azbill usage-details --format flatten,pretty -S XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX --start 2020-06-01 --end 2020-06-30 -o usage.jsonl
2020/07/06 23:43:59 Loading auth-dev token in /Users/yaegashi/.azbill/auth_dev.json
2020/07/06 23:43:59 Requesting with consumption.UsageDetailsClient
2020/07/06 23:43:59    scope: "subscriptions/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
2020/07/06 23:43:59   filter: "properties/usageStart eq '2020-06-01' and properties/usageEnd eq '2020-06-30'"
2020/07/06 23:44:11 Writing to "usage.jsonl" in json,flatten,pretty
.                                                       138 records
2020/07/06 23:44:11 Done 138 records in 16.893702ms, 8168.724653 records/sec
```

For Enterprise Agreement accounts: Export usage details of billing period 202006 for billing account XXXXXXXX into usage.csv in CSV format:

```console
$ azbill usage-details --format csv -A XXXXXXXX -P 202006 -o usage.csv
2020/07/06 23:45:15 Loading auth-dev token in /Users/yaegashi/.azbill/auth_dev.json
2020/07/06 23:45:15 Requesting with consumption.UsageDetailsClient
2020/07/06 23:45:15    scope: "providers/Microsoft.Billing/billingAccounts/XXXXXXXX/providers/Microsoft.Billing/billingPeriods/202006"
2020/07/06 23:45:15   filter: ""
2020/07/06 23:45:34 Writing to "usage.csv" in csv
..................................................     5000 records
..................................................    10000 records
..................................................    15000 records
....................                                  17001 records
2020/07/06 23:46:38 Done 17001 records in 1m4.325743007s, 264.295431 records/sec
```

## Development

azbill utilizes [Azure REST API](https://docs.microsoft.com/en-us/rest/api/azure/)
and the API version not yet available in [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go)
to acquire all the usage details including the reservation
in [the consumption API](https://docs.microsoft.com/en-us/rest/api/consumption/).

The CLI command hierarchy is built with [spf13/cobra](https://github.com/spf13/cobra)
and [yaegashi/cobra-comder](https://github.com/yaegashi/cobra-cmder).
