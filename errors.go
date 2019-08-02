package shortener

import "errors"

var ErrLinkNotFound = errors.New("Link not found")
var ErrIDMissing = errors.New("Missing link ID parameter")
var ErrDataMissing = errors.New("Missing POST body data")
var ErrIDExists = errors.New("Given Id already exists")
var ErrInvalidCredentials = errors.New("Invalid Credentials!")
