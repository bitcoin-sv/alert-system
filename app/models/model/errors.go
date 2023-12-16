package model

import "errors"

// ErrCannotConvertToIDs is the error when the conversion fails from interface into type IDs
var ErrCannotConvertToIDs = errors.New("cannot convert value to type IDs")

// ErrMissingDatastore missing datastore from model
var ErrMissingDatastore = errors.New("datastore is missing from model, cannot save")
