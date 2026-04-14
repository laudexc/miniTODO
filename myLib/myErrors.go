package myLib

import "errors"

var ErrBookNotFound = errors.New("book not found")
var ErrBookAlreadyExists = errors.New("book already exists")
