package dot

import (
	"errors"
	"io/fs"
	"testing"
)

func TestChoseBackupVersion (t *testing.T,){
	
}







type MockDirEntry struct {
	name     string
	isDir    bool
	fileInfo fs.FileInfo
}

func NewMockDirEntry(name string, isDir bool, fileInfo fs.FileInfo) *MockDirEntry {
	return &MockDirEntry{
		name:     name,
		isDir:    isDir,
		fileInfo: fileInfo,
	}
}

func (m *MockDirEntry) Info() (fs.FileInfo, error) {
	if m.fileInfo == nil {
		return nil, errors.New("file info not set")
	}
	return m.fileInfo, nil
}

func (m *MockDirEntry) IsDir() bool {
	return m.isDir
}

func (m *MockDirEntry) Name() string {
	return m.name
}

func (m *MockDirEntry) Type() fs.FileMode {
	if m.fileInfo != nil {
		return m.fileInfo.Mode()
	}
	return 0
}
