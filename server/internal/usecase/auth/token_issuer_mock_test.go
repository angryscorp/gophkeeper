package auth

type stubIssuer struct {
	issueFn func(userID, deviceID string) (string, error)
}

func (s stubIssuer) IssueAccess(userID, deviceID string) (string, error) {
	if s.issueFn != nil {
		return s.issueFn(userID, deviceID)
	}
	return "token", nil
}
