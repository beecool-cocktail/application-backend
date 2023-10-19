package domain

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	CodeSuccess                                = "00000"
	CodeRequestDecodeFailed                    = "P0001"
	CodeCanNotSpecifyHttpAction                = "P0002"
	CodeParameterIllegal                       = "P0003"
	CodePermissionDenied                       = "P0004"
	CodeCommandNotFound                        = "C0001"
	CodeItemDoesNotBelongToUser                = "N0001"
	CodeUserAlreadyExist                       = "A0001"
	CodeUserNotFound                           = "A0002"
	CodePasswordNotMatch                       = "A0003"
	CodeTokenExpired                           = "A0004"
	CodeCocktailNotFound                       = "B0001"
	CodeCocktailDraftIsMaximum                 = "B0002"
	CodeCocktailNotFinished                    = "B0003"
	CodeFavoriteCocktailListNotOpenToThePublic = "B0004"
	CodeFileTypeIllegal                        = "F0001"
	CodeFileSizeIllegal                        = "F0002"
	CodeInternalError                          = "S0001"
	CodeResponseEncodedFailed                  = "S0002"
	CodeRedisLockNotObtained                   = "S0003"
)

var (
	ErrRequestDecodeFailed                      = errors.New("request decode failed")
	ErrResponseEncodedFailed                    = errors.New("response encoded failed")
	ErrCanNotSpecifyHttpAction                  = errors.New("can't specify action through request parameter")
	ErrParameterIllegal                         = errors.New("parameter illegal")
	ErrPermissionDenied                         = errors.New("permission denied")
	ErrItemDoesNotBelongToUser                  = errors.New("item doesn't belong to user")
	ErrUserAlreadyExist                         = errors.New("user already exist")
	ErrUserNotFound                             = errors.New("user not found")
	ErrPasswordNotMatch                         = errors.New("password not match")
	ErrCocktailNotFound                         = errors.New("cocktail not found")
	ErrorCocktailDraftIsMaximum                 = errors.New("cocktail draft is maximum")
	ErrorCocktailNotFinished                    = errors.New("cocktail not finished")
	ErrorFavoriteCocktailListNotOpenToThePublic = errors.New("favorite cocktail not open to the public")
	ErrTokenExpired                             = errors.New("token expired")
	ErrCommandNotFound                          = errors.New("command expire or not exist")
	ErrCodeFileTypeIllegal                      = errors.New("illegal file type")
	ErrCodeFileSizeIllegal                      = errors.New("illegal file size")
	ErrFilePathIllegal                          = errors.New("illegal file path")
	ErrRedisLockNotObtained                     = errors.New("lock not obtained")
	ErrInternalError														= errors.New("internal error")
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
	case ErrCanNotSpecifyHttpAction:
		return CodeParameterIllegal
	case ErrParameterIllegal:
		return CodeParameterIllegal
	case ErrPermissionDenied:
		return CodePermissionDenied
	case ErrCommandNotFound:
		return CodeCommandNotFound
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
	case ErrorCocktailDraftIsMaximum:
		return CodeCocktailDraftIsMaximum
	case ErrorCocktailNotFinished:
		return CodeCocktailNotFinished
	case ErrorFavoriteCocktailListNotOpenToThePublic:
		return CodeFavoriteCocktailListNotOpenToThePublic
	case ErrCodeFileTypeIllegal:
		return CodeFileTypeIllegal
	case ErrCodeFileSizeIllegal:
		return CodeFileSizeIllegal
	case ErrRedisLockNotObtained:
		return CodeRedisLockNotObtained

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
	case ErrCanNotSpecifyHttpAction:
		return http.StatusBadRequest
	case ErrParameterIllegal:
		return http.StatusBadRequest
	case ErrCommandNotFound:
		return http.StatusNotFound
	case ErrPermissionDenied:
		return http.StatusForbidden
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
	case ErrorCocktailDraftIsMaximum:
		return http.StatusBadRequest
	case ErrorCocktailNotFinished:
		return http.StatusBadRequest
	case ErrorFavoriteCocktailListNotOpenToThePublic:
		return http.StatusUnauthorized
	case ErrCodeFileTypeIllegal:
		return http.StatusForbidden
	case ErrCodeFileSizeIllegal:
		return http.StatusForbidden
	case ErrRedisLockNotObtained:
		return http.StatusOK

	default:
		return http.StatusInternalServerError
	}
}
