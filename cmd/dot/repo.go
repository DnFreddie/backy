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

func getHeadUrl(url string) (string, error) {
	re := regexp.MustCompile(`github.com/(.*)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		repoPath := matches[1]
		fmt.Println("Repository Path:", repoPath)

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repoPath)
		fmt.Println("API URL:", apiURL)

		client := &http.Client{}
		resp, err := client.Get(apiURL)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d\n", resp.StatusCode)
			return "", fmt.Errorf("received status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return "", err
		}

		var details struct {
			DefaultBranch string `json:"default_branch"`
			RepoName      string `json:"name"`
		}

		err = json.Unmarshal(body, &details)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return "", err
		}

		fmt.Println("Default Branch:", details.DefaultBranch)

		zipURL := fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.zip", repoPath, details.DefaultBranch)

		return zipURL, nil
	} else {
		fmt.Println("No match found in the URL")
		return "", fmt.Errorf("no match found in the URL")
	}
}

func downloadGit(zipURL string, repoName string) error {
	file, err := os.Create(repoName + ".zip") // Use a .zip extension

	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	client := &http.Client{}
	r, err := client.Get(zipURL)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", r.StatusCode)
		return fmt.Errorf("received status code %d", r.StatusCode)
	}

	_, err = io.Copy(file, r.Body)
	if err != nil {
		fmt.Println("Error copying response body to file:", err)
		return err
	}
	err = utils.UnzipSource(file.Name(), repoName)

	if err != nil {
		err = os.Remove(file.Name())

		if err != nil {

			return fmt.Errorf("Failed while downlowading %s remove the file yourself ", file.Name())
		}
		return err
	}

	return nil
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

func (d *Dotfile) ignore(toIgnore []string) {
	d.ignored = false

	for _, pattern := range toIgnore {
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
