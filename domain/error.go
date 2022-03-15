package domain

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	CodeSuccess                 = "00000"
	CodeRequestDecodeFailed     = "P0001"
	CodeParameterIllegal        = "P0003"
	CodeItemDoesNotBelongToUser = "N0001"
	CodeUserAlreadyExist        = "A0001"
	CodeUserNotFound            = "A0002"
	CodePasswordNotMatch        = "A0003"
	CodeTokenExpired            = "A0004"
	CodeCocktailNotFound        = "B0001"
	CodeFileTypeIllegal         = "F0001"
	CodeFileSizeIllegal         = "F0002"
	CodeInternalError           = "S0001"
	CodeResponseEncodedFailed   = "S0002"
)

var (
	ErrRequestDecodeFailed     = errors.New("request decode failed")
	ErrResponseEncodedFailed   = errors.New("response encoded failed")
	ErrParameterIllegal        = errors.New("parameter illegal")
	ErrItemDoesNotBelongToUser = errors.New("item doesn't belong to user")
	ErrUserAlreadyExist        = errors.New("user already exist")
	ErrUserNotFound            = errors.New("user not found")
	ErrPasswordNotMatch        = errors.New("password not match")
	ErrCocktailNotFound        = errors.New("cocktail not found")
	ErrTokenExpired            = errors.New("token expired")
	ErrCodeFileTypeIllegal     = errors.New("illegal file type")
	ErrCodeFileSizeIllegal     = errors.New("illegal file size")
)

func GetErrorCode(err error) string {
	if err == nil {
		return CodeSuccess
	}

	logrus.Error(err)

	switch err {
	case ErrRequestDecodeFailed:
		return CodeRequestDecodeFailed
	case ErrResponseEncodedFailed:
		return CodeResponseEncodedFailed
	case ErrItemDoesNotBelongToUser:
		return CodeItemDoesNotBelongToUser
	case ErrParameterIllegal:
		return CodeParameterIllegal
	case ErrUserAlreadyExist:
		return CodeUserAlreadyExist
	case ErrUserNotFound:
		return CodeUserNotFound
	case ErrPasswordNotMatch:
		return CodePasswordNotMatch
	case ErrTokenExpired:
		return CodeTokenExpired
	case ErrCocktailNotFound:
		return CodeCocktailNotFound
	case ErrCodeFileTypeIllegal:
		return CodeFileTypeIllegal
	case ErrCodeFileSizeIllegal:
		return CodeFileSizeIllegal

	default:
		return CodeInternalError
	}
}

func GetStatusCode(err error) int {

	logrus.Error(err)

	switch err {
	case ErrRequestDecodeFailed:
		return http.StatusBadRequest
	case ErrResponseEncodedFailed:
		return http.StatusInternalServerError
	case ErrItemDoesNotBelongToUser:
		return http.StatusForbidden
	case ErrParameterIllegal:
		return http.StatusBadRequest
	case ErrUserAlreadyExist:
		return http.StatusBadRequest
	case ErrUserNotFound:
		return http.StatusNotFound
	case ErrPasswordNotMatch:
		return http.StatusUnauthorized
	case ErrTokenExpired:
		return http.StatusUnauthorized
	case ErrCocktailNotFound:
		return http.StatusNotFound
	case ErrCodeFileTypeIllegal:
		return http.StatusForbidden
	case ErrCodeFileSizeIllegal:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
