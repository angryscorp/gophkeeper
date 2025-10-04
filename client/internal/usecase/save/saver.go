package save

import (
	"context"
	"gophkeeper/client/internal/domain"
	"time"
)

const (
	ctxTimeout = 5 * time.Second
)

type CredentialsRepository interface {
	Save(ctx context.Context, credentials domain.Credentials) error
}

type UserInfoSaver struct {
	credentialsRepo CredentialsRepository
}

func New(credentialsRepo CredentialsRepository) *UserInfoSaver {
	return &UserInfoSaver{
		credentialsRepo: credentialsRepo,
	}
}

var _ domain.UserInfoSaver = (*UserInfoSaver)(nil)

func (u UserInfoSaver) SaveCredentials(credentials domain.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	err := u.credentialsRepo.Save(ctx, credentials)
	if err != nil {
		return err
	}

	return nil
}

func (u UserInfoSaver) SaveBankCard(bankCard domain.BankCard) error {
	//TODO implement me
	return nil
}

func (u UserInfoSaver) SaveUserBinaryData(userBinaryData domain.UserBinaryData) error {
	//TODO implement me
	return nil
}

func (u UserInfoSaver) SaveUserTextData(userTextData domain.UserTextData) error {
	//TODO implement me
	return nil
}
