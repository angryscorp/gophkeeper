package menu

import "gophkeeper/client/internal/domain"

// Environment provides the dependencies and factory functions
// required by the menu and its sub-models. It wires together
// application services (auth, sync, data storage, help) with
// the TUI screens.
type Environment struct {
	RegFactory   func(username, password string) error
	LoginFactory func(username, password string) error
	DataSaver    domain.UserInfoSaver
	SyncFactory  func() error
	DataFactory  func() ([]domain.Record, error)
	HelpFactory  func() string
}
