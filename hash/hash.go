package hash

import (
)

type FileChanged struct {
	AbPath     string
	Hash       []byte
	WasChanged bool

}
