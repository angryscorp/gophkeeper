package domain

type UserInfoSaver interface {
	SaveBankCard(bankCard BankCard) error
	SaveCredentials(credentials Credentials) error
	SaveUserBinaryData(userBinaryData UserBinaryData) error
	SaveUserTextData(userTextData UserTextData) error
}
