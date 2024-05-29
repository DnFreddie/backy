package trip

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sync"
)

type FileChanged struct {
	AbPath     string
	Hash       []byte
	WasChanged bool
}

func HashHash(b []byte, fPath string) FileChanged {
	h := sha256.New()
	hashed := h.Sum(b)

	testChanged := FileChanged{
		AbPath:     fPath,
		Hash:       hashed,
		WasChanged: false,
	}

	return testChanged
}

func readFilesAsync(filePath string) (FileChanged,error){

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist:", filePath)
		return FileChanged{},nil
	}

	content, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Can't read the file permissions issues", err)
		return FileChanged{}, err
	}

	fChanged := HashHash(content, filePath)


	return fChanged ,nil
}


type Directory struct {
    Path         string
    Subdirectories map[string]*Directory
    Files        []FileChanged
    mu           sync.Mutex 
}

func TripCommand(fPath string) (*Directory, error) {
    dir := &Directory{
        Path:          fPath,
        Subdirectories: make(map[string]*Directory),
    }

    fds, err := os.ReadDir(fPath)
    if err != nil {
        return nil, err
    }

    var wg sync.WaitGroup 
    for _, fd := range fds {
        jPath := path.Join(fPath, fd.Name())

        if fd.IsDir() {
            wg.Add(1)
            fdCopy := fd
            go func(jPath string, fdCopy os.DirEntry) {
                defer wg.Done() 
                subDir, err := TripCommand(jPath)
                if err != nil {
                    fmt.Printf("Error in directory %s: %v\n", jPath, err)
                    return
                }
                dir.mu.Lock()
                defer dir.mu.Unlock()
                dir.Subdirectories[fdCopy.Name()] = subDir
            }(jPath, fdCopy)
        } else {
            has, err := readFilesAsync(jPath)
            if err != nil {
                fmt.Printf("Error reading file %s: %v\n", jPath, err)
                continue
            }
            dir.mu.Lock()
            dir.Files = append(dir.Files, has)
            dir.mu.Unlock()
        }
    }

    wg.Wait()

    return dir, nil
}

func printDirectory(dir *Directory) {
	fmt.Println("-------------------")
	fmt.Println("this is the dir ",dir.Path)
	fmt.Println("-------------------")
    for _, file := range dir.Files {
        fmt.Println( file.AbPath)
    }
    
    for _, subDir := range dir.Subdirectories {
        printDirectory(subDir)
    }
}
