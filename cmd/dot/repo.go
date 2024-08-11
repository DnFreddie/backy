package dot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DnFreddie/backy/utils"
)

func gitClone(url string) (string, error) {

	done := make(chan bool)

	utils.WaitingScreen(done, "Cloning")
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

	done <- true

	return pathToRepo, nil

}

func (r *Repo) Clone(url string) error{
	done := make(chan bool)

	utils.WaitingScreen(done, "Cloning")
	if err := r.getHeadUrl(url); err != nil {
		return fmt.Errorf("failed to get head URL: %w", err)
	}

	if err := r.downloadRepo(); err != nil {
		return fmt.Errorf("failed to download repo %s: %w", r.RepoName, err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	zipFile := r.RepoName + ".zip"
	repoPath, err := utils.UnzipSource(zipFile, pwd)
	if err != nil {
		return fmt.Errorf("failed to unzip %s: %w", r.RepoName, err)
	}

defer func() {
    os.Remove(zipFile)
    done <- true
}()


	 r.AbsP,err = utils.MakeAbsolute(repoPath)
	
if err != nil {

	return fmt.Errorf("%v doesn't exist: %w", repoPath, err)

}
	return nil
}

type Repo struct {
	DefaultBranch string `json:"default_branch"`
	Url           string
	zipUrl        string
	RepoName      string `json:"name"`
	AbsP          string
}

func (r *Repo) getHeadUrl(url string) error {
	re := regexp.MustCompile(`github.com/(.*)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		repoPath := matches[1]

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repoPath)
		fmt.Println("Repository:", repoPath)

		client := &http.Client{}
		resp, err := client.Get(apiURL)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d\n", resp.StatusCode)
			return fmt.Errorf("received status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return err
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return err
		}

		fmt.Println("Default Branch:", r.DefaultBranch)

		zipURL := fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.zip", repoPath, r.DefaultBranch)
		r.zipUrl = zipURL
		return nil
	} else {
		fmt.Println("No match found in the URL")
		return fmt.Errorf("no match found in the URL")
	}
}

// Downloads the file named archive .zip
func (r *Repo) downloadRepo() error {

	client := &http.Client{}
	res, err := client.Get(r.zipUrl)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", res.StatusCode)
		return fmt.Errorf("received status code %d", res.StatusCode)
	}

	file, err := os.Create(r.RepoName + ".zip")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)

	if err != nil {

		return utils.HandleFileErr("copying data to file", err, file)

	}

	return nil
}

func readIgnore() []string {
	var ignored []string

	ignored = append(ignored, IGNORE)
	ignored = append(ignored, ".git")

	_, err := os.Stat(IGNORE)
	if os.IsNotExist(err) {
		fmt.Println("No git ignore")
		return ignored
	}

	c, err := os.ReadFile(IGNORE)
	if err != nil {
		fmt.Println("Can't read the file", err)
		return ignored
	}

	sc := string(c)
	lines := strings.Split(sc, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			ignored = append(ignored, trimmedLine)
		}
	}

	return ignored
}
func (d *Dotfile) ignore(toIgnore *[]string) {
	d.ignored = false

	for _, pattern := range *toIgnore {
		if match, _ := filepath.Match(pattern, d.Location.Name()); match {
			d.ignored = true
			break
		}
	}
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

	if strings.Contains(str, "git@") {
		return true
	}
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
