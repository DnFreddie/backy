package dot

import (
	"fmt"
	"os"
	"path"

	"github.com/DnFreddie/backy/utils"
)

func CopyTemp(dotfiels []Dotfile) {
	userConf, err := utils.GetUser("Desktop")

	if err != nil {
		return
	}

	cDir, err := utils.ScanDir(userConf)

	if err != nil {
		return
	}

	fmt.Println(cDir)
	dst := "/tmp/test_data"
	err = os.MkdirAll(dst, 0700)

	if err != nil {
		return

	}
	var ErrMessage []error

	var errMessages []error
	for _, dot := range cDir {
		fullPath := path.Join(userConf, dot.Name())
		fmt.Println(fullPath)
		if dot.IsDir() {
			err = utils.CopyDir(fullPath, path.Join(dst, dot.Name()))
		} else {
			err = utils.CopyFile(fullPath, path.Join(dst, dot.Name()))
		}
		if err != nil {
			errMessages = append(errMessages, fmt.Errorf("failed to copy %s: %v", dot.Name(), err))
		}

		fmt.Println(ErrMessage)
	}
}
