package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	//"github.com/DnFreddie/backy/utils"
)

func (r *Repo) revert() {

	for _, d := range *r.Dots {
		if d.New {

		}
		fmt.Println("Moving", d.Symlink)
		fmt.Println("Moving", d.Absolute)
	}

}

func chooseBackupVersion(options []os.DirEntry) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Choose the backup version\n")
		for i, dir := range options {
			prettyName, err := time.Parse("20060102150405", dir.Name())
			if err != nil {
				log.Fatal("You must have modified one of the directories. Don't do that", err)
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

		return fmt.Sprintf("%v", options[choice-1].Name()), nil
	}
}

func processReversion(chosenPath string) (Repo, error) {
	r := &Repo{}
	csvPath := path.Join(chosenPath, "backy_schema.josn")
	fmt.Println(csvPath)

	_, err := os.Stat(csvPath)

	if os.IsNotExist(err) {

		log.Fatal("The schema for reversion doesn't exist")
	}

	f, err := os.Open(csvPath)
	if err != nil {

	}
	defer f.Close()

	bytes, err := io.ReadAll(f)

	json.Unmarshal(bytes, r)

	return *r, nil
}
