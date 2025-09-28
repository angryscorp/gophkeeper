package record

type recordRoute int

const (
	routeTypeSelection recordRoute = iota
	routeBankCardForm
	routeCredentialsForm
	routeTextDataForm
	routeBinaryDataForm
)
