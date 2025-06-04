package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret, iss, aud string
}

func NewJWTAuthenticator(secret, iss, aud string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		iss:    iss,
		aud:    aud,
	}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secret), nil
	}, jwt.WithAudience(j.aud), jwt.WithIssuer(j.iss), jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
}
