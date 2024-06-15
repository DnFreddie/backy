package dot

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/DnFreddie/backy/utils"
)

const (
	IGNORE    = ".gitignore"
	BACK_CONF = "back_conf"
	REVERT_CSV = "schema.csv"
)

func createTempBack(source string, backupDir string, csvF *csv.Writer) error {

	_, err := os.Stat(source)

	if os.IsNotExist(err) {
		fmt.Println("No git ignore ")
		fmt.Println(err)
		return nil
	}

	fmt.Println("why thsi source doens't work", source)
	dest := path.Join(backupDir, path.Base(source))
	err = os.Rename(source, dest)

	if err != nil {
		return err
	}
	data := [][]string{
		{source, dest},
	}
	err = csvF.WriteAll(data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

// Returns the path to the repo
func gitClone(url string) (string, error) {
	cmd := exec.Command("bash", "-c", "git clone "+url)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	re := regexp.MustCompile(`[^/]+$`)

	match := re.FindString(url)

	pwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	if strings.HasSuffix(match, ".git") {
		match = strings.TrimSuffix(match, ".git")
	}
	pathToRepo := path.Join(pwd, match)

	fmt.Println(pathToRepo)
	return pathToRepo, nil

}

func DotCommand(repo string) error {
	var isURL bool
	var dest string

	if strings.Contains(repo, "git@") {
		isURL = true
	} else {
		_, err := url.ParseRequestURI(repo)
		if err == nil {
			isURL = true
		}
	}

	if isURL {
		clonedDest, err := gitClone(repo)
		if err != nil {
			return err
		}
		dest = clonedDest
	} else {
		dest = repo
	}

	dirPaths, err := GetPaths(dest)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dirStructs := Isexe(dirPaths)

	err = CreateSymlink(dirStructs, dest)
	if err != nil {
		return err
	}

	return nil
}

func CreateSymlink(dotfiles []Dotfile, source string) error {

	targetPath, err := utils.GetUser("Desktop")
	if err != nil {

		return err
	}

	nowT := time.Now().Format("20060102150405")

	backupDir, err := utils.Checkdir(path.Join(BACK_CONF, nowT), false)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(path.Join(backupDir, REVERT_CSV))
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(f)

	for _, f := range dotfiles {

		if f.IsEx {
			symlinkPath := f.Location.Name()
			sourceAbs := path.Join(source, symlinkPath)
			dest := path.Join(targetPath, symlinkPath)

			//Checks weateter the file exist and moves it to backup dir 
			err := createTempBack(dest, backupDir, writer)
			if err != nil {
				fmt.Println("Error creating temporary backup:", err)
				return err
			}

			err = os.Symlink(sourceAbs, dest)
			if err != nil {
				fmt.Println("Failed to create symlink:", err)
				return err
			}
		}
	}

	return nil

}

func GetPaths(gitPath string) ([]fs.DirEntry, error) {
	dirs, err := os.ReadDir(gitPath)
	//fmt.Println(dirs)
	if err != nil {
		fmt.Println("Can't list this dir probably permissions issue ", err)
		return nil, err
	}

	toIgnore, err := readIgnore()
	if err != nil {
		fmt.Println("Can't read git ignore: ", err)
		return nil, err
	}

	var paths []fs.DirEntry
	for _, dir := range dirs {
		if !shouldIgnore(dir.Name(), toIgnore) {
			paths = append(paths, dir)
		}
	}
	//fmt.Printf("\nthih are the paths {%v}", paths)
	return paths, nil
}

func shouldIgnore(fileName string, toIgnore []string) bool {
	for _, pattern := range toIgnore {
		if match, _ := filepath.Match(pattern, fileName); match {
			//fmt.Printf("its a match {%v}  {%v} filename", pattern, fileName)
			return true

		}
	}
	return false
}

func readIgnore() ([]string, error) {

	_, err := os.Stat(IGNORE)
	if os.IsNotExist(err) {
		fmt.Println("No git ignore ")
		return nil, nil
	}

	c, err := os.ReadFile(IGNORE)
	if err != nil {
		fmt.Println("Can't read the file", err)
		return nil, err
	}

	sc := string(c)

	ignored := strings.Split(sc, "\n")
	ignored = append(ignored, ".git")
	ignored = append(ignored, IGNORE)

	//fmt.Println("this are ignored", ignored)
	return ignored, nil
}
