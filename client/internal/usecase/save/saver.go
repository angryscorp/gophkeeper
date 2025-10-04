package save

import "gophkeeper/client/internal/domain"

type UserInfoSaver struct {
}

var _ domain.UserInfoSaver = (*UserInfoSaver)(nil)

func (u UserInfoSaver) SaveCredentials(credentials domain.Credentials) error {
	//TODO implement me
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
