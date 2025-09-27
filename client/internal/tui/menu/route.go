package menu

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
