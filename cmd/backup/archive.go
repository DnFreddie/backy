package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ZipDir(sourceDirs []string, zipDest string) error {
	zipFileWriter, err := os.Create(zipDest)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFileWriter.Close()

	zipWriter := zip.NewWriter(zipFileWriter)
	defer zipWriter.Close()

	for _, sourceDir := range sourceDirs {
		baseDir := filepath.Base(sourceDir)

		err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(sourceDir, filePath)
			if err != nil {
				return err
			}
			zipFilePath := filepath.Join(baseDir, relPath)

			zipFileHeader, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			zipFileHeader.Name = zipFilePath

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
			return fmt.Errorf("failed to zip directory %s: %v", sourceDir, err)
		}

		fmt.Printf("Directory %s added successfully \n", sourceDir)
	}

	fmt.Printf("All directories zipped successfully to %s\n", zipDest)
	return nil
}
