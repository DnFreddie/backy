package dot

import (
	"bufio"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

func RevertBackups() error {
	confDir, err := utils.Checkdir(BACK_CONF, false)
	if err != nil {
		return err
	}

	fmt.Println(confDir)
	options, err := os.ReadDir(confDir)
	if err != nil {
		return err
	}

	fmt.Println(len(options))
	if len(options) < 1 {
		fmt.Println("No backups to revert")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	var revDir fs.DirEntry

	for {
		fmt.Print("Choose the backup version\n")
		for i, dir := range options {
			fmt.Printf("%d: %s\n", i+1, dir.Name())
		}

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		choice, err := strconv.Atoi(text)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice, please choose a valid option.")
			continue
		}

		revDir = options[choice-1]
		break
	}

	fmt.Printf("You chose: %s\n", revDir.Name())

	return nil
}
