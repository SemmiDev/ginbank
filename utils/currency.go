package utils

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	IDR = "IDR"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	//case USD, EUR, CAD, IDR:
	case IDR:
		return true
	}
	return false
}
