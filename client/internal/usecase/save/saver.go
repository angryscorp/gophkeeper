package save

import (
	"context"
	"encoding/json"
	"gophkeeper/client/internal/domain"
	"time"
)

const (
	ctxTimeout = 5 * time.Second
)

type UserInfoSaver struct {
	repo      Repository
	encryptor func([]byte) ([]byte, error)
}

func New(
	repo Repository,
	encryptor func([]byte) ([]byte, error),
) *UserInfoSaver {
	return &UserInfoSaver{
		repo:      repo,
		encryptor: encryptor,
	}
}

var _ domain.UserInfoSaver = (*UserInfoSaver)(nil)

func (u UserInfoSaver) SaveCredentials(credentials domain.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	payload, err := u.encryptor(data)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindCredentials, credentials.ID, payload)
}

func (u UserInfoSaver) SaveBankCard(bankCard domain.BankCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(bankCard)
	if err != nil {
		return err
	}

	payload, err := u.encryptor(data)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindBankCard, bankCard.ID, payload)
}

func (u UserInfoSaver) SaveUserBinaryData(userBinaryData domain.UserBinaryData) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(userBinaryData)
	if err != nil {
		return err
	}

	payload, err := u.encryptor(data)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindBinaryData, userBinaryData.ID, payload)
}

func (u UserInfoSaver) SaveUserTextData(userTextData domain.UserTextData) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(userTextData)
	if err != nil {
		return err
	}

	payload, err := u.encryptor(data)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindTextData, userTextData.ID, payload)
}
