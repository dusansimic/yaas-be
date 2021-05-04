package auth_test

import (
	"testing"

	"github.com/dusansimic/yaas/services/auth"
	"github.com/stretchr/testify/require"
)

func TestJWTService_GenerateToken_ok(t *testing.T) {
	s := auth.JWTAuthService()

	_, err := s.GenerateToken("name", 1)

	require.NoError(t, err)
}

func TestJWTService_ValidateToken_ok(t *testing.T) {
	s := auth.JWTAuthService()

	tk, err := s.GenerateToken("name", 1)

	require.NoError(t, err)

	_, err = s.ValidateToken(tk)

	require.NoError(t, err)
}

func TestJWTService_ValidateToken_invalidToken(t *testing.T) {
	s := auth.JWTAuthService()

	_, err := s.ValidateToken("")

	require.Error(t, err)
}

const invalidToken = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5hbWUiLCJpZCI6MX0.ZFag95HhMfHkZ1_BLbxabD-Qj7i_DddwvDHSeqHY-ZKj0W9lYSk6PUWvVoPgxPYmTOdnVGd3QSFFllm4YYrBRg"

func TestJWTService_ValidateToken_invalidMethod(t *testing.T) {
	s := auth.JWTAuthService()

	_, err := s.ValidateToken(invalidToken)

	require.Error(t, err)
	require.EqualError(t, err, "invalid token ES256")
}
