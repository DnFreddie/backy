package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ZipDir(sourceDir, zipDest string) error {
	zipFileWriter, err := os.Create(zipDest)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFileWriter.Close()

	zipWriter := zip.NewWriter(zipFileWriter)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		zipFileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		zipFileHeader.Name, err = filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			zipFileHeader.Name += "/"
		}

		fileWriter, err := zipWriter.CreateHeader(zipFileHeader)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileToZip, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer fileToZip.Close()

			_, err = io.Copy(fileWriter, fileToZip)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to zip directory: %v", err)
	}

	fmt.Printf("Directory %s zipped successfully to %s\n", sourceDir, zipDest)
	return nil
}
