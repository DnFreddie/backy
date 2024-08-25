package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/DnFreddie/backy/utils"
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
		slog.Error("Faieled to get the url", "err", err)
		return err
	}

	if err := r.downloadRepo(); err != nil {
		slog.Error("failed to download", "repo", r.RepoName, "err", err)
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current directory: %w", err)
	}

	zipFile := r.RepoName + ".zip"
	repoPath, err := utils.UnzipSource(zipFile, pwd)
	if err != nil {
		slog.Error("failed to unzip", "repo", r.RepoName, "err", err)
		return err
	}

	defer func() {
		os.Remove(zipFile)
		done <- true
	}()

	r.Absolute, err = utils.MakeAbsolute(repoPath)
	if err != nil {

		slog.Error("Can't make absoulute", "path", repoPath, "err", err)
		return err

	}
	return nil
}

func (r *Repo) ReadLocal(localPath string) error {
	absoluteP, err := utils.MakeAbsolute(localPath)
	if err != nil {
		slog.Error("Doesn't exist:", "path", path.Base(localPath))
		return err
	}
	r.Absolute = absoluteP

	HEAD := filepath.Join(r.Absolute, ".git/HEAD")
	CONFIG := filepath.Join(r.Absolute, ".git/config")
	reBranch := regexp.MustCompile(`refs/heads/(\w+)`)
	reUrl := regexp.MustCompile(`url = (.+\.git)$`)

	// Attempt to read the HEAD file
	headFile, err := os.OpenFile(HEAD, os.O_RDONLY, 0)
	if err != nil {
		slog.Warn("Failed to open HEAD file:", "error", err)
	} else {
		defer headFile.Close()
		reader := bufio.NewReader(headFile)
		r.DefaultBranch, err = extractMatch(reader, reBranch)
		if err != nil {
			slog.Warn("Failed to extract default branch:", "error", err)
		}
	}

	configFile, err := os.OpenFile(CONFIG, os.O_RDONLY, 0)
	if err != nil {
		slog.Warn("Failed to open config file:", "error", err)
	} else {
		defer configFile.Close()
		configReader := bufio.NewReader(configFile)
		r.Url, err = extractMatch(configReader, reUrl)
		if err != nil {
			slog.Warn("Failed to extract URL:", "error", err)
		}
	}

	return nil
}

func (r *Repo) Link(force bool) {
	if r.Dots == nil {
		slog.Info("No Files to link", "repo", r.RepoName)
		return
	}

	target := path.Join(os.Getenv("HOME"), CONFIG)

	r.createBackup()

	for i := range *r.Dots {
		d := &(*r.Dots)[i]
		d.IsExe()
		if force || d.Executable {
			if err := d.createSymlink(target); err != nil {
				slog.Warn("Failed to symlink", "target", target)
			}
		}
	}

}
func (r *Repo) PrintReport() {

	if len(*r.Dots) == 0 {
		fmt.Println(utils.Cyan + "No dotfiles to raport " + utils.Reset)
		return
	}

	fmt.Println(utils.Cyan + "Generating report..." + utils.Reset)
	var fCounter int
	for _, i := range *r.Dots {
		if i.Failed != "" {
			fmt.Println(utils.Red + "Error Report:" + utils.Reset)
			fmt.Printf("%sName: %s %s\n", utils.Green, i.Location, utils.Reset)
			fmt.Printf("%sFailure: %v%s\n\n", utils.Red, i.Failed, utils.Reset)

			fCounter++
		}

	}
	success := len(*r.Dots) - fCounter
	fmt.Printf("%sRepo: %s%s\n", utils.Blue, r.RepoName, utils.Reset)
	fmt.Printf("%sSucceed: %d%s\n", utils.Green, success, utils.Reset)
	fmt.Printf("%sFailed: %v%s\n\n", utils.Red, fCounter, utils.Reset)

}
func (r *Repo) getHeadUrl(url string) error {
	re := regexp.MustCompile(`github.com/(.*)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		repoPath := matches[1]

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repoPath)
		slog.Info("Repository:", "repoPath", repoPath)

		client := &http.Client{}
		resp, err := client.Get(apiURL)
		if err != nil {
			slog.Error("Error making GET request:", "error", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			slog.Error("Error: received status code", "statusCode", resp.StatusCode)
			return fmt.Errorf("received status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body:", "error", err)
			return err
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			slog.Error("Error unmarshalling JSON:", "error", err)
			return err
		}

		slog.Info("Default Branch:", "defaultBranch", r.DefaultBranch)

		zipURL := fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.zip", repoPath, r.DefaultBranch)
		r.zipUrl = zipURL
		r.Url = url
		return nil
	} else {
		slog.Error("No match found in the URL", "url", url)
		return fmt.Errorf("no match found in the URL")
	}
}

// Downloads the file named archive .zip
func (r *Repo) downloadRepo() error {
	client := &http.Client{}
	res, err := client.Get(r.zipUrl)
	if err != nil {
		slog.Error("Error making GET request:", "error", err)
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Error: received status code", "statusCode", res.StatusCode)
		return fmt.Errorf("error: received status code %d", res.StatusCode)
	}

	file, err := os.Create(r.RepoName + ".zip")
	if err != nil {
		slog.Error("Error creating file:", "error", err)
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.Error("Error closing file:", "error", cerr)
			err = fmt.Errorf("error closing file: %w", cerr)
		}
	}()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		slog.Error("Error during copying data to file:", "error", err)
		if cleanupErr := utils.HandleFileErr("copying data to file", err, file); cleanupErr != nil {
			slog.Error("Error during cleanup:", "error", cleanupErr)
			return fmt.Errorf("error during cleanup: %w", cleanupErr)
		}
		return err
	}

	slog.Info("Repository downloaded successfully:", "repoName", r.RepoName)
	return nil
}

func (r *Repo) readIgnore() {
	ignored := []string{IGNORE, ".git"}

	gitIgnore := filepath.Join(r.Absolute, IGNORE)
	if _, err := os.Stat(gitIgnore); os.IsNotExist(err) {
		slog.Info("No git ignore found, using default ignored files")
		r.GitIgnore = ignored
		return
	} else if err != nil {
		slog.Error("Error checking git ignore:", "error", err)
		r.GitIgnore = ignored
		return
	}

	c, err := os.OpenFile(gitIgnore, os.O_RDONLY, 0)
	if err != nil {
		slog.Warn("Can't read git ignore, skipping:", "error", err)
		r.GitIgnore = ignored
		return
	}
	defer c.Close()

	parsedIgnored, err := parseReadIgnore(io.Reader(c))
	if err != nil {
		slog.Warn("Error parsing git ignore:", "error", err)
		r.GitIgnore = ignored
		return
	}

	r.GitIgnore = append(ignored, parsedIgnored...)
	slog.Info("Git ignore processed successfully")
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
		slog.Error("Can't list this dir probably permissions issue", "error", err)
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
	backupDir := path.Join(DOTS, nowT)
	backupPath, err := utils.Checkdir(backupDir, false)
	if err != nil {
		log.Fatal("Failed to create backup")
	}

	r.BackupLocation = backupPath

}

func (r *Repo) SaveSchema() error {
	if r.BackupLocation == "" {

		return fmt.Errorf("Can't find the backup location")
	}

	jsonData, jsonErr := json.MarshalIndent(*r, "", "  ")
	if jsonErr != nil {
		slog.Error("Failed to marshal data:", "error", jsonErr)
		return jsonErr
	}

	schemaDest := path.Join(r.BackupLocation, utils.SCHEMA_JSON)
	file, writeErr := os.Create(schemaDest)
	if writeErr != nil {
		slog.Error("Failed to create output.json:", "error", writeErr)
		return writeErr
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.Error("Error closing file:", "error", cerr)
		}
	}()

	if _, writeErr := file.Write(jsonData); writeErr != nil {
		slog.Error("Failed to write JSON data to file:", "error", writeErr)
		return writeErr
	}

	slog.Info("Schema saved successfully to", "schemaDest", schemaDest)
	return nil
}
