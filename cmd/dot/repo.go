package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/DnFreddie/backy/utils"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Yellow = "\033[33m"
)

type Repo struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	DefaultBranch string `json:"default_branch"`
	Url           string `json:"repo_url"`
	zipUrl        string `json:"-"`
	RepoName      string `json:"name"`
	Absolute      string `json:"-"`
	// TODO! add this when the parser will be complete
	GitIgnore      []string   `gorm:"-" json:"-"`
	BackupLocation string     `json:"backup_location"`
	Dots           *[]Dotfile `gorm:"-"`
	RepoId         string     `gorm:"foreignKey:RepoID" json:"repo_id"`
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

		fmt.Println("errr", err)
		return err

	}
	return nil
}

func (r *Repo) ReadLocal(localPath string) error {
	absoluteP, err := utils.MakeAbsolute(localPath)
	if err != nil {
		return fmt.Errorf("%v doesn't exist: %w", path.Base(localPath), err)
	}
	r.Absolute = absoluteP

	HEAD := filepath.Join(r.Absolute, ".git/HEAD")
	CONFIG := filepath.Join(r.Absolute, ".git/config")
	reBranch := regexp.MustCompile(`refs/heads/(\w+)`)
	reUrl := regexp.MustCompile(`url = (.+\.git)$`)

	headFile, err := os.OpenFile(HEAD, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open HEAD file: %w", err)
	}
	defer headFile.Close()

	reader := bufio.NewReader(headFile)
	r.DefaultBranch, err = extractMatch(reader, reBranch)
	if err != nil {
		return fmt.Errorf("failed to extract default branch: %w", err)
	}

	configFile, err := os.OpenFile(CONFIG, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	configReader := bufio.NewReader(configFile)
	r.Url, err = extractMatch(configReader, reUrl)
	if err != nil {
		return fmt.Errorf("failed to extract URL: %w", err)
	}

	return nil
}

func (r *Repo) Link(force bool) {
	if r.Dots == nil {
		fmt.Println("No Files to link")
		return
	}


	target, err := utils.GetUser(TARGET)
	if err != nil {
		log.Fatal(err)
	}

	r.createBackup()


	for i := range *r.Dots {
		d := &(*r.Dots)[i] 
		d.IsExe()
		if force || d.Executable {
			if err := d.createSymlink(target); err != nil {
				fmt.Println(err)
			}
		}
	}

}
func (r *Repo) PrintRaport() {

	if len(*r.Dots) == 0 {
		fmt.Println(Cyan + "No dotfiles to raport " + Reset)
		return
	}

	fmt.Println(Cyan + "Generating report..." + Reset)
	var fCounter int
	for _, i := range *r.Dots {
		if i.Failed != nil {
			fmt.Println(Red + "Error Report:" + Reset)
			fmt.Printf("Name: %s\n", Green+i.Location+Reset)
			fmt.Printf("Failure: %s\n\n", Red+fmt.Sprintf("%v", *i.Failed)+Reset)
			fCounter++
		}

	}
	success := len(*r.Dots) - fCounter
	fmt.Printf("Repo: %s\n", Blue+fmt.Sprintf("%v", r.RepoName)+Reset)
	fmt.Printf(Green+"Succeed: %d\n"+Reset, success)
	fmt.Printf(Red+"Failed: %s\n\n", fmt.Sprintf("%v", fCounter)+Reset)
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
	ignored := []string{IGNORE, ".git"}

	gitIgnore := filepath.Join(r.Absolute, IGNORE)
	if _, err := os.Stat(gitIgnore); os.IsNotExist(err) {
		fmt.Println("No git ignore found")
		r.GitIgnore = ignored
		return
	} else if err != nil {
		fmt.Println("Error checking git ignore:", err)
		r.GitIgnore = ignored
		return
	}

	c, err := os.OpenFile(gitIgnore, os.O_RDONLY, 0)

	if err != nil {
		fmt.Println("Can't read git ignore, skipping:", err)
		r.GitIgnore = ignored
		return
	}

	parsedIgnored, err := parseReadIgnore(io.Reader(c))
	if err != nil {
		fmt.Println("Error parsing git ignore:", err)
		r.GitIgnore = ignored
		return
	}

	r.GitIgnore = append(ignored, parsedIgnored...)
}
func parseReadIgnore(re io.Reader) ([]string, error) {
	var ignored []string
	scanner := bufio.NewScanner(re)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			ignored = append(ignored, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ignore file: %w", err)
	}

	return ignored, nil
}

func extractMatch(reader io.Reader, re *regexp.Regexp) (string, error) {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		line := s.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			captureGroup := matches[1]
			return captureGroup, nil
		}
	}

	if err := s.Err(); err != nil {
		return "", err
	}

	return "", nil
}
func (r *Repo) getDots() error {

	dirs, err := os.ReadDir(r.Absolute)
	if err != nil {
		fmt.Println("Can't list this dir probably permissions issue ", err)
		return err

	}
	var dotfiels []Dotfile

	for _, d := range dirs {

		dot := Dotfile{
			Location: d.Name(),
			Absolute: path.Join(r.Absolute, d.Name()),
			Repo:     r,
			RepoID:   &r.RepoId,
		}
		dotfiels = append(dotfiels, dot)

	}

	r.readIgnore()
	for _, dot := range dotfiels {
		dot.ignore(&r.GitIgnore)
	}
	r.Dots = &dotfiels

	return nil
}

func (r *Repo) createBackup() {
	nowT := time.Now().Format("20060102150405")
	backupDir := path.Join(BACK_CONF, nowT)
	backupPath, err := utils.Checkdir(backupDir, false)
	if err != nil {
		log.Fatal("Failed to create backup")
	}

	r.BackupLocation = backupPath

}

func (r *Repo) saveRepoSchema() {
	if r.BackupLocation == "" {
		log.Fatal("Can't find the backup location")
	}
	
	// for _,i := range *r.Dots{
	// 	fmt.Println("This is the symlink")
	// 	fmt.Println(i.Symlink)
	// }

		
	jsonData, jsonErr := json.MarshalIndent(*r, "", "  ")
	if jsonErr != nil {
		log.Println("Failed to marshal data:", jsonErr)
		return 
	}

	schemaDest := path.Join(r.BackupLocation, utils.SCHEMA_JSON)
	file, writeErr := os.Create(schemaDest)
	if writeErr != nil {
		log.Println("Failed to create output.json:", writeErr)
		return 
	}
	defer file.Close()

	if _, writeErr := file.Write(jsonData); writeErr != nil {
		log.Println("Failed to write JSON data to file:", writeErr)
		return 
	}

	log.Println("Schema saved successfully to", schemaDest)
}

