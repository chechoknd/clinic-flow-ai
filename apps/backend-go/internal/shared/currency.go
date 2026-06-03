package shared

import "strings"

type CurrencyMetadata struct {
	Code              string `json:"code"`
	Symbol            string `json:"symbol"`
	Locale            string `json:"locale"`
	DecimalDigits     int    `json:"decimal_digits"`
	ThousandSeparator string `json:"thousand_separator"`
	DecimalSeparator  string `json:"decimal_separator"`
	SymbolPosition    string `json:"symbol_position"`
}

type CountryCurrency struct {
	CountryCode string
	Currency    CurrencyMetadata
}

var supportedCurrencies = map[string]CountryCurrency{
	"CO": {CountryCode: "CO", Currency: CurrencyMetadata{Code: "COP", Symbol: "$", Locale: "es-CO", DecimalDigits: 0, ThousandSeparator: ".", DecimalSeparator: ",", SymbolPosition: "before"}},
	"PE": {CountryCode: "PE", Currency: CurrencyMetadata{Code: "PEN", Symbol: "S/", Locale: "es-PE", DecimalDigits: 2, ThousandSeparator: ",", DecimalSeparator: ".", SymbolPosition: "before"}},
	"AR": {CountryCode: "AR", Currency: CurrencyMetadata{Code: "ARS", Symbol: "$", Locale: "es-AR", DecimalDigits: 2, ThousandSeparator: ".", DecimalSeparator: ",", SymbolPosition: "before"}},
	"CL": {CountryCode: "CL", Currency: CurrencyMetadata{Code: "CLP", Symbol: "$", Locale: "es-CL", DecimalDigits: 0, ThousandSeparator: ".", DecimalSeparator: ",", SymbolPosition: "before"}},
}

func NormalizeCountryCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func NormalizeCurrencyCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func CurrencyForCountry(countryCode string) (CurrencyMetadata, bool) {
	item, ok := supportedCurrencies[NormalizeCountryCode(countryCode)]
	return item.Currency, ok
}

func CurrencyByCode(currencyCode string) (CurrencyMetadata, bool) {
	code := NormalizeCurrencyCode(currencyCode)
	for _, item := range supportedCurrencies {
		if item.Currency.Code == code {
			return item.Currency, true
		}
	}
	return CurrencyMetadata{}, false
}

func IsValidCountryCurrency(countryCode, currencyCode string) bool {
	item, ok := supportedCurrencies[NormalizeCountryCode(countryCode)]
	return ok && item.Currency.Code == NormalizeCurrencyCode(currencyCode)
}

func SupportedCountryCurrencies() []CountryCurrency {
	return []CountryCurrency{
		supportedCurrencies["CO"],
		supportedCurrencies["PE"],
		supportedCurrencies["AR"],
		supportedCurrencies["CL"],
	}
}
