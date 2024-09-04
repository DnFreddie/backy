package hash
import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
)

type FileNode struct {
	Name     string      `json:"name"`
	IsDir    bool        `json:"is_dir"`
	Children []*FileNode `json:"children,omitempty"`
	Hash     *string     `json:"hash,omitempty"`
}

func WalkDir(root string) (*FileNode, error) {
	rootNode := &FileNode{
		Name:  filepath.Base(root),
		IsDir: true,
	}

	nodeMap := map[string]*FileNode{
		".": rootNode,
	}

	err := fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		node := &FileNode{
			Name:  d.Name(),
			IsDir: d.IsDir(),
		}

		parent := nodeMap[filepath.Dir(path)]
		parent.Children = append(parent.Children, node)

		if d.IsDir() {
			nodeMap[path] = node
		} else {
			fullPath := filepath.Join(root, path)
			hash, err := HashHash(fullPath)
			if err != nil {
				return err
			}
			node.Hash = &hash
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func HashHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

