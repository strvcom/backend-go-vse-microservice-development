package authenticator

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret string
}

func NewJWTAuthenticator(secret string) JWTAuthenticator {
	return JWTAuthenticator{secret: secret}
}

func (a JWTAuthenticator) VerifyToken(token string) (map[string]interface{}, error) {
	var claims jwtClaims
	if _, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	}); err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"role":   claims.Role,
		"userID": claims.UserID,
	}
	return result, nil
}

type jwtClaims struct {
	jwt.RegisteredClaims

	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
