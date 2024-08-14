package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DnFreddie/backy/utils"
)

type Repo struct {
	ID             uint       `gorm:"primaryKey;autoIncrement"`
	DefaultBranch string `json:"default_branch"`
	Url           string
	zipUrl        string
	RepoName      string `json:"name"`
	Absolute      string
	GitIgnore     []string `gorm:"-"`
	BackupLocation  string 
	Dots *[]Dotfile `gorm:"-"`
	RepoId  string  `gorm:"foreignKey:RepoID"`

} 

func (r *Repo) Clone(url string) error {
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

	r.Absolute, err = utils.MakeAbsolute(repoPath)
	if err != nil {

		fmt.Println("errr",  err)
		return err

	}
	return nil
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
		r.Url = url
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
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received status code %d", res.StatusCode)
	}

	file, err := os.Create(r.RepoName + ".zip")

	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("error closing file: %w", cerr)
		}
	}()

	_, err = io.Copy(file, res.Body)
	if err != nil {

		if cleanupErr := utils.HandleFileErr("copying data to file", err, file); cleanupErr != nil {
			return fmt.Errorf("error during cleanup: %w", cleanupErr)
		}
		return err 
	}

	return nil
}


func (r *Repo) readIgnore() {
	var ignored []string

	ignored = append(ignored, IGNORE)
	ignored = append(ignored, ".git")
	gitIgnore := filepath.Join(r.Absolute, IGNORE)
	_, err := os.Stat(gitIgnore)
	if os.IsNotExist(err) {
		fmt.Println("No git ignore found")
		r.GitIgnore = ignored
		return
	}

	c, err := os.ReadFile(gitIgnore)
	if err != nil {
		fmt.Println("Can't read git ignore skipping")
		r.GitIgnore = ignored
	}

	sc := string(c)
	lines := strings.Split(sc, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		if trimmedLine != "" {
			ignored = append(ignored, trimmedLine)
		}
	}
	r.GitIgnore = ignored

}

func (r *Repo) GetInfo(localPath string) error {
	absouluteP, err := utils.MakeAbsolute(localPath)
	if err != nil {
		return fmt.Errorf("%v doesn't exist: %w", path.Base(localPath), err)
	}
	r.Absolute = absouluteP
	HEAD := filepath.Join(r.Absolute, "HEAD")
	CONFIG := filepath.Join(r.Absolute, "config")
	reBranch := regexp.MustCompile(`refs/heads/(\w+)`)
	reUrl := regexp.MustCompile(`url = (.+.git$)`)

	r.DefaultBranch,err = extractMatch(HEAD, reBranch)
	r.Url,err  = extractMatch(CONFIG, reUrl)

	if err != nil {
		fmt.Println(err)
		return nil

	}
	return nil
}
func extractMatch(filePath string, re *regexp.Regexp) (string,error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	defer file.Close()
	if err != nil {
		
	return "",fmt.Errorf("Faield to read the %s", path.Base(filePath))	
	}
	buffReader := bufio.NewScanner(file)
	for buffReader.Scan() {
		line := buffReader.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			captureGroup := matches[1]
			return captureGroup,nil
		}
		continue

	}
	return "",err
}
func (r *Repo) getDots() ( error) {

	dirs, err := os.ReadDir(r.Absolute)
	if err != nil {
		fmt.Println("Can't list this dir probably permissions issue ", err)
		return  err

	}
	var dotfiels []Dotfile

	for _, d := range dirs {

		dot := Dotfile{
			Location: d.Name(),
			AbPath:   path.Join(r.Absolute, d.Name()),
			Repo:     r,
			RepoID: r.RepoId,
		}
		dotfiels = append(dotfiels, dot)

	}

	r.readIgnore()
	for _, dot := range dotfiels {
		dot.ignore(&r.GitIgnore)
	}
r.Dots = &dotfiels

	return  nil
}
