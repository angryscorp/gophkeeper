package domain

type UserInfoType string

const (
	UserInfoTypeBankCard       UserInfoType = "BankCard"
	UserInfoTypeCredentials    UserInfoType = "Credentials"
	UserInfoTypeUserBinaryData UserInfoType = "UserBinaryData"
	UserInfoTypeUserTextData   UserInfoType = "UserTextData"
)
