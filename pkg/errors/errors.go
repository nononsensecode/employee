package errors

type Code int
type ErrorType string

const (
	Unknown  ErrorType = "unknown"
	NotFound ErrorType = "not-found"
	Auth     ErrorType = "authentication"
	Input    ErrorType = "input"
)

type ApiError interface {
	error
	Code() Code
	ErrorType() ErrorType
}

type defaultError struct {
	code      Code
	err       error
	errorType ErrorType
}

func (d defaultError) Code() Code {
	return d.code
}

func (d defaultError) Error() string {
	return d.err.Error()
}

func (d defaultError) ErrorType() ErrorType {
	return d.errorType
}

func NewNotFoundError(code Code, err error) ApiError {
	return defaultError{
		code:      code,
		errorType: NotFound,
		err:       err,
	}
}

func NewUnknownError(code Code, err error) ApiError {
	return defaultError{
		code:      code,
		errorType: Unknown,
		err:       err,
	}
}

func NewInputError(code Code, err error) ApiError {
	return defaultError{
		code:      code,
		errorType: Input,
		err:       err,
	}
}

func NewAuthError(code Code, err error) ApiError {
	return defaultError{
		code:      code,
		errorType: Auth,
		err:       err,
	}
}
