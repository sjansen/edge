package inspector

import "errors"

var ErrInvalidSignature = errors.New("invalid function signature")

var ErrTypeErrorsFound = errors.New("type errors found")
