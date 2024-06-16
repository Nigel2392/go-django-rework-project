package pages

import "github.com/Nigel2392/django/core/errs"

const (
	ErrPathLengthTooShort errs.Error = "path is too short"
	ErrPathLengthExceeded errs.Error = "path length exceeded"
	ErrInvalidPathLength  errs.Error = "invalid path length"
	ErrTooLittleAncestors errs.Error = "too little ancestors provided"
	ErrTooManyAncestors   errs.Error = "too many ancestors provided"
	ErrPageIsRoot         errs.Error = "page is root"
	ErrInvalidScanType    errs.Error = "invalid scan type"
)
