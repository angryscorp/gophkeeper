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
	repo Repository
}

func New(repo Repository) *UserInfoSaver {
	return &UserInfoSaver{
		repo: repo,
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

	return u.repo.Save(ctx, domain.UserDataKindCredentials, credentials.ID, data)
}

func (u UserInfoSaver) SaveBankCard(bankCard domain.BankCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(bankCard)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindBankCard, bankCard.ID, data)
}

func (u UserInfoSaver) SaveUserBinaryData(userBinaryData domain.UserBinaryData) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(userBinaryData)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindBinaryData, userBinaryData.ID, data)
}

func (u UserInfoSaver) SaveUserTextData(userTextData domain.UserTextData) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	data, err := json.Marshal(userTextData)
	if err != nil {
		return err
	}

	return u.repo.Save(ctx, domain.UserDataKindTextData, userTextData.ID, data)
}
