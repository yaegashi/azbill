package mapconv

import (
	"reflect"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

var (
	TimeType    = reflect.TypeOf(date.Time{})
	DecimalType = reflect.TypeOf(decimal.Decimal{})
	UUIDType    = reflect.TypeOf(uuid.UUID{})
)
