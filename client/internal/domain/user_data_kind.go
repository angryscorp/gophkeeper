package domain

type UserDataKind = int64

const (
	UserDataKindBankCard UserDataKind = iota + 1
	UserDataKindCredentials
	UserDataKindTextData
	UserDataKindBinaryData
)
