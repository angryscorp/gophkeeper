package record

// recordRoute identifies which screen of the record-creation
// flow is currently active: the type selection menu or one
// of the specific data entry forms.
type recordRoute int

const (
	routeTypeSelection recordRoute = iota
	routeBankCardForm
	routeCredentialsForm
	routeTextDataForm
	routeBinaryDataForm
)
