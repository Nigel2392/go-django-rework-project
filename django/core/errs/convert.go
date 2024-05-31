package errs

import (
	"errors"
	"fmt"
)

//func convert(err any, default_ any) error {
//	if err == nil && default_ == nil {
//		return nil
//	}
//
//	if err == nil {
//		return Convert(default_, nil)
//	}
//
//	switch e := err.(type) {
//	case error:
//		return e
//	case string:
//		return errors.New(e)
//	default:
//		return fmt.Errorf("%v", e)
//	}
//}

func Convert(fmtOrErr any, default_ any, args ...any) error {
	if fmtOrErr == nil && default_ == nil {
		return nil
	}

	if fmtOrErr == nil {
		return Convert(default_, nil, args...)
	}

	switch e := fmtOrErr.(type) {
	case error:
		return e
	case string:
		if len(args) > 0 {
			return fmt.Errorf(e, args...)
		} else {
			return errors.New(e)
		}
	default:
		if len(args) > 0 {
			return errors.New(
				fmt.Sprint(append([]any{e}, args...)),
			)
		} else {
			return fmt.Errorf("%v", e)
		}
	}
}
