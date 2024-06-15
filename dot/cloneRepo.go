package dot

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func gitClone(url string) (string, error) {

	go func() {
		for {
			for _, r := range `-\|/` {
				fmt.Printf("\rCloning %c", r)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

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

	return pathToRepo, nil

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

	return ignored, nil
}
func shouldIgnore(fileName string, toIgnore []string) bool {
	for _, pattern := range toIgnore {
		if match, _ := filepath.Match(pattern, fileName); match {
			return true

		}
	}
	return false
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
