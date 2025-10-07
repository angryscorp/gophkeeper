package auth

// TokenIssuer defines the ability to create signed access tokens
// (e.g. JWT) for a given user and device. The returned string is
// the serialized token that can be passed to clients.
type TokenIssuer interface {
	IssueAccess(userID, deviceID string) (string, error)
}
