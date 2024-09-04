package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipSource(source, destination string) (string, error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	destination, err = filepath.Abs(destination)
	if err != nil {
		return "", err
	}

	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return "", err
		}
	}

	first := reader.File[0].Name
	return first, nil
}

func unzipFile(f *zip.File, destination string) error {

	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}


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
