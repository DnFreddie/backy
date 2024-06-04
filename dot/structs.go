package dot

import "io/fs"

type Dotfile struct {
	Location fs.DirEntry
	IsEx     bool
	Symlink  string
	BaseP    string
}
