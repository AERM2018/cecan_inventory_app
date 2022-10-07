package customreqvalidations

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/exp/maps"
)

func ValidateUser(mapObject interface{}, omitedFields ...string) error {
	return validation.Validate(mapObject,
		validation.Map(
			validation.Key("role_id", validation.Required, is.UUIDv4),
			validation.Key("name", validation.Required),
			validation.Key("surname", validation.Required),
			validation.Key("email", validation.Required, is.Email),
			validation.Key("password", validation.Required, validation.Length(8, 50)),
		).AllowExtraKeys())
}

func ValidateMedicine(mapObject interface{}, omitedFields ...string) error {
	return validation.Validate(mapObject,
		validation.Map(
			validation.Key("key", validation.Required, validation.Length(9, 9)),
			validation.Key("name", validation.Required),
		).AllowExtraKeys())
}

func ValidatePharmacyStock(mapObject interface{}, omitedFields ...string) error {
	rules := map[string]*validation.KeyRules{
		"medicine_key": validation.Key("medicine_key", validation.Required.Error("medicine key is required")),
		"lot_number":   validation.Key("lot_number", validation.Required),
		"pieces":       validation.Key("pieces", validation.Required.Error("must be more than 0 pieces")),
		"expires_at":   validation.Key("expires_at", validation.Required, validation.Date(time.RFC3339Nano).Min(time.Now()).Error("the date must be greater than today")),
	}
	for _, key := range omitedFields {
		delete(rules, key)
	}
	return validation.Validate(mapObject,
		validation.Map(maps.Values(rules)...).AllowExtraKeys())
}

func ValidatePrescription(mapObject interface{}, ommitFields ...string) error {
	return validation.Validate(mapObject,
		validation.Map(
			validation.Key("patient_name", validation.Required),
			validation.Key("instructions", validation.Required),
			validation.Key("medicines", validation.Length(0, 100)),
		).AllowExtraKeys())
}

func ValidateStorehouseUtility(mapObject interface{}, ommitFields ...string) error {
	return validation.Validate(mapObject,
		validation.Map(
			validation.Key("key", validation.Required),
			validation.Key("generic_name", validation.Required),
			validation.Key("storehouse_utility_category_id", validation.Required),
			validation.Key("storehouse_utility_presentation_id", validation.Required),
			validation.Key("storehouse_utility_unit_id", validation.Required),
			validation.Key("quantity_per_unit", validation.Required.Error("The quantity per unit must be grater than 0")),
		).AllowExtraKeys())
}

func ValidateStorehouseStock(mapObject interface{}, omitedFields ...string) error {
	rules := map[string]*validation.KeyRules{
		"storehouse_utility_key": validation.Key("storehouse_utility_key", validation.Required.Error("storehouse utility key is required")),
		"quantity_presentation":  validation.Key("quantity_presentation", validation.Required.Error("must be more than 0")),
		"lot_number":             validation.Key("lot_number", validation.Required),
		"catalog_number":         validation.Key("catalog_number", validation.Required),
		"expires_at":             validation.Key("expires_at", validation.Required, validation.Date(time.RFC3339Nano).Min(time.Now()).Error("the date must be greater than today")),
	}
	for _, key := range omitedFields {
		delete(rules, key)
	}
	return validation.Validate(mapObject,
		validation.Map(maps.Values(rules)...).AllowExtraKeys())
}
