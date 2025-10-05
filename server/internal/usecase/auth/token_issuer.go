package auth

type TokenIssuer interface {
	IssueAccess(userID, deviceID string) (string, error)
}
