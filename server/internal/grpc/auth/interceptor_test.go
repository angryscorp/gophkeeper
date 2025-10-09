package auth

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/server/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapError(t *testing.T) {
	tests := []struct {
		in   error
		want codes.Code
	}{
		{domain.ErrUsernameTaken, codes.AlreadyExists},
		{domain.ErrUsernameNotFound, codes.NotFound},
		{domain.ErrChallengeNotFound, codes.FailedPrecondition},
		{domain.ErrChallengeFailed, codes.Unauthenticated},
		{context.Canceled, codes.Canceled},
		{nil, codes.OK},
		{errors.New("random"), codes.Unknown},
	}

	for _, tt := range tests {
		got := mapError(tt.in)
		if tt.in == nil && got != nil {
			t.Errorf("expected nil, got %v", got)
			continue
		}
		if tt.in != nil {
			st, _ := status.FromError(got)
			if st.Code() != tt.want {
				t.Errorf("for %v expected %v, got %v", tt.in, tt.want, st.Code())
			}
		}
	}
}
