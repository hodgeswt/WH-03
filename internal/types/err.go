package types

import "errors"

var (
    ErrDataTooLarge = errors.New("ErrDataTooLarge")
    ErrInvalidChars = errors.New("ErrInvalidChars")
)
