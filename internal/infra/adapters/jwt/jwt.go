package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/valueobject"
	"github.com/andreis3/customers-ms/internal/infra/configs"
)

type JWT struct {
	secret []byte
	expiry time.Duration
}

func NewJWT(conf *configs.Configs) *JWT {
	return &JWT{
		secret: []byte(conf.JWTSecret),
		expiry: conf.JWTExpiry,
	}
}

func (j *JWT) CreateToken(customer entity.Customer) (*valueobject.TokenClaims, *errors.Error) {
	claims := j.createJWTMapClaims(customer)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return nil, errors.ErrorCreateToken(err)
	}

	tokenClaims := j.createTokenClaims(customer, tokenString)

	return &tokenClaims, nil
}

func (j *JWT) ValidateToken(tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	tokenClaims, err := j.parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return tokenClaims, nil
}

func (j *JWT) RefreshToken(tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	tokenClaims, err := j.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, toMapClaims(tokenClaims))

	newToken, errSig := token.SignedString(j.secret)

	if errSig != nil {
		return nil, errors.ErrorRefreshToken(errSig)
	}

	tokenClaims.Token = newToken

	return tokenClaims, nil
}

func (j *JWT) DecodeToken(tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	return j.parseToken(tokenString)
}

func (j *JWT) parseToken(tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrorInvalidTokenAlgorithmError()
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, errors.ErrorValidateToken(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.ErrorValidateToken(err)
	}

	// Safe access with checks
	customerID, ok := claims["customer_id"].(float64)
	if !ok {
		return nil, errors.ErrorValidateTokenMessage("invalid or missing customer_id")
	}
	fullName, ok := claims["full_name"].(string)
	if !ok {
		return nil, errors.ErrorValidateTokenMessage("invalid or missing full_name")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.ErrorValidateTokenMessage("invalid or missing email")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.ErrorValidateTokenMessage("invalid or missing exp")
	}

	return &valueobject.TokenClaims{
		CustomerID: int64(customerID),
		FullName:   fullName,
		Email:      email,
		ExpiresAt:  time.Unix(int64(exp), 0),
	}, nil
}

func (j *JWT) createTokenClaims(customer entity.Customer, tokenString string) valueobject.TokenClaims {
	return valueobject.TokenClaims{
		CustomerID: customer.ID,
		FullName:   customer.FullName(),
		Email:      customer.Email.String(),
		Token:      tokenString,
		ExpiresAt:  time.Now().Add(j.expiry),
	}
}

func (j *JWT) createJWTMapClaims(customer entity.Customer) jwt.MapClaims {
	now := time.Now()
	return jwt.MapClaims{
		"customer_id": customer.ID,
		"full_name":   customer.FullName(),
		"email":       customer.Email.String(),
		"exp":         now.Add(j.expiry).Unix(),
		"iat":         now.Unix(),
		"nbf":         now.Unix(),
	}
}

func toMapClaims(tokenClaims *valueobject.TokenClaims) jwt.MapClaims {
	return jwt.MapClaims{
		"customer_id": tokenClaims.CustomerID,
		"full_name":   tokenClaims.FullName,
		"email":       tokenClaims.Email,
		"exp":         tokenClaims.ExpiresAt.Unix(),
	}
}
