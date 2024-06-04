package trip

import (
	"encoding/json"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"os"
)

type ConfigPath struct {
	Fpath string `json:"fpath"`
}

func createConfig(scanPath string) (bool, error) {

	confP, err := utils.Checkdir("scan_paths.json")
	var existed []ConfigPath

	err = utils.ReadJson(confP, &existed)
	if err != nil {
		fmt.Println("Can't read json")
		return true, err
	}

	var exist bool
	for _, p := range existed {

		if p.Fpath == scanPath {
			exist = true
			break
		}
	}

	if !exist {
		existed = append(existed, ConfigPath{Fpath: scanPath})
		bytes, err := json.Marshal(&existed)
		if err != nil {
			return false, err
		}
		os.WriteFile(confP, bytes, 0666)
	}

	return exist, nil
}
