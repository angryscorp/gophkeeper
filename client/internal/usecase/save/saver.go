package save

import (
	"context"
	"encoding/json"
	"time"

	"gophkeeper/client/internal/domain"
)

const (
	ctxTimeout = 5 * time.Second
)

// UserInfoSaver encrypts and persists user data (cards, credentials, etc.)
// into the underlying Repository.
type UserInfoSaver struct {
	repo      Repository
	encryptor func([]byte) ([]byte, error)
}

// New creates a new UserInfoSaver with the given repository and encryptor.
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

// SaveCredentials marshals, encrypts, and saves a Credentials record.
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

// SaveBankCard marshals, encrypts, and saves a BankCard record.
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

// SaveUserBinaryData marshals, encrypts, and saves a UserBinaryData record.
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

// SaveUserTextData marshals, encrypts, and saves a UserTextData record.
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
