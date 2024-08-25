package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ChooseBackupVersion(options []os.DirEntry) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose the backup version")

		for i := len(options) - 1; i >= 0; i-- {
			dir := options[i]
			prettyName, err := time.Parse("20060102150405", dir.Name())

			if err != nil {
				log.Fatal(UserError{Err: err, FPath: dir.Name()})
			}
			fmt.Printf("%d: %s%s%s\n", i+1, Cyan, prettyName.Format("January 2, 2006 15:04:05"), Reset)
		}

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		choice, err := strconv.Atoi(text)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice, please choose a valid option.")
			continue
		}

		return options[len(options)-choice].Name(), nil
	}
}
