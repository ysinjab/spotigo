package auth

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type key int

const (
	secret           = "123"
	TokenInfoKey key = iota
)

func parseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			return claims, nil
		} else {
			return nil, fmt.Errorf("invalid toke")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, fmt.Errorf("that's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, fmt.Errorf("token is either expired or not active yet")
		} else {
			return nil, fmt.Errorf("couldn't handle this token %s", err)
		}
	} else {
		return nil, fmt.Errorf("couldn't handle this token %s", err)
	}
}

func userClaimFromToken(claims jwt.MapClaims) interface{} {
	return claims["sub"]
}

func ValidateToken(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
	newCtx := context.WithValue(ctx, TokenInfoKey, tokenInfo)
	return newCtx, nil
}

func AuthFunc(ctx context.Context) (context.Context, error) {
	newCtx, err := ValidateToken(ctx)
	if err != nil {
		return nil, err
	}
	return newCtx, nil
}
