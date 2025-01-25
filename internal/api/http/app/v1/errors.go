package v1

// Errors
const (
	UnknownErrorCode    = 0
	UnknownErrorMessage = "unknown error"

	ClientAlreadyExistsCode                 = 1001
	ClientAlreadyExistsMessage              = "user already exists"
	ClientNotFoundCode                      = 1002
	ClientNotFoundMessage                   = "user not found"
	ClientRefreshTokenCookieNotFoundCode    = 1003
	ClientRefreshTokenCookieNotFoundMessage = "user refresh token cookie not found"
	ClientRefreshTokenExpiredCode           = 1004
	ClientRefreshTokenExpiredMessage        = "user refresh token expired"
)

type ErrorCode int
type ErrorMessage string

type ErrorStruct struct {
	ErrorCode    `json:"error_code"`
	ErrorMessage `json:"error_message"`
} // @name ErrorStruct

type ValidationErrorStruct struct {
	ErrorCode    int               `json:"error_code"`
	ErrorMessage string            `json:"error_message"`
	Errors       []ValidationError `json:"validation_errors"`
}

type ValidationError struct {
	FieldKey     string `json:"field_key"`
	ErrorMessage string `json:"error_message"`
}

func getErrorStruct(code ErrorCode) *ErrorStruct {
	errorStruct := &ErrorStruct{
		ErrorCode:    UnknownErrorCode,
		ErrorMessage: UnknownErrorMessage,
	}

	switch code {
	case ClientAlreadyExistsCode:
		errorStruct.ErrorCode = ClientAlreadyExistsCode
		errorStruct.ErrorMessage = ClientAlreadyExistsMessage
	case ClientNotFoundCode:
		errorStruct.ErrorCode = ClientNotFoundCode
		errorStruct.ErrorMessage = ClientNotFoundMessage
	case ClientRefreshTokenCookieNotFoundCode:
		errorStruct.ErrorCode = ClientRefreshTokenCookieNotFoundCode
		errorStruct.ErrorMessage = ClientRefreshTokenCookieNotFoundMessage
	case ClientRefreshTokenExpiredCode:
		errorStruct.ErrorCode = ClientRefreshTokenExpiredCode
		errorStruct.ErrorMessage = ClientRefreshTokenExpiredMessage
	}

	return errorStruct
}
