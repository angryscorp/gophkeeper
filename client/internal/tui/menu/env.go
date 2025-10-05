package menu

import "gophkeeper/client/internal/domain"

type Environment struct {
	RegFactory   func(username, password string) error
	LoginFactory func(username, password string) error
	DataSaver    domain.UserInfoSaver
	SyncFactory  func() error
	DataFactory  func() ([]domain.Record, error)
	HelpFactory  func() string
}
