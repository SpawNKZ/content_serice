package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRouting            = errors.New("inconsistent mapping between route and handler")
	ErrInconsistentIDs       = errors.New("inconsistent IDs")
	ErrAlreadyExists         = errors.New("already exists")
	ErrNotFound              = errors.New("object not found")
	ErrDB                    = errors.New("db error")
	ErrBadRequest            = errors.New("bad request")
	ErrContentIsNotRemovable = errors.New("content is not removable")
	ErrLimit                 = fmt.Errorf("limit should be greater than zero and less than")
	ErrOffset                = fmt.Errorf("offset should be greater than or equal to zero and less than")
	ErrInvalidLocale         = errors.New("invalid language code")
	ErrInconsistentLocales   = errors.New("inconsistent locales")
	ErrTranslationExists     = errors.New("duplicate translation")
	ErrParsing               = errors.New("failed to parse")
	ErrInvalidOrganizationID = errors.New("invalid Organization ID")
	ErrInvalidCityID         = errors.New("invalid City ID")
	ErrInvalidGradeNumber    = errors.New("invalid Grade Number")
	ErrInvalidSchoolID       = errors.New("invalid School ID")
	ErrInvalidGradeID        = errors.New("invalid Grade ID")
	ErrInvalidClassID        = errors.New("invalid Class ID")
	ErrInvalidSchoolYear     = errors.New("invalid School Year")
	ErrInvalidLetter         = errors.New("invalid Letter")
	ErrNoResponseFromSubject = errors.New("could not get response from the subject service, wrong response format")
)

func ErrorToHttpCode(err error) int {
	switch err {
	case ErrBadRouting, ErrInconsistentIDs, ErrLimit, ErrOffset, ErrInvalidLocale, ErrInconsistentLocales,
		ErrTranslationExists, ErrParsing, ErrInvalidOrganizationID, ErrInvalidCityID, ErrInvalidGradeNumber,
		ErrInvalidSchoolID, ErrInvalidGradeID, ErrInvalidSchoolYear:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
