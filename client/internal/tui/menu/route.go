package menu

// route identifies the active screen in the TUI.
// It is used by the menu model to switch between
// different sub-models (auth, sync, data, etc.).
type route int

const (
	routeMenu route = iota
	routeRegister
	routeAuth
	routeSync
	routeData
	routeNewItem
	routeHelp
	routeQuit
)
