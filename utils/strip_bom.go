package utils

import (
	"io"
	"os"

	"github.com/dimchansky/utfbom"
)

// The UTF-8 BOM byte sequence
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// StripBOM cleans the BOM, writes to a temp file, and returns the path to the temp file
func StripBOM(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Wrap the file reader with the BOM stripper
	// This creates an io.Reader that automatically skips the BOM if it exists
	bomFreeReader := utfbom.SkipOnly(file)

	// Create a temporary file to hold the BOM-free content
	// This is necessary because gpandas.Read_csv only takes a path.
	err = os.Mkdir("temp", 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	tmpFile, err := os.CreateTemp("temp", "clean_csv_*.csv")
	if err != nil {
		return "", err
	}
	defer func(tmpFile *os.File) {
		err := tmpFile.Close()
		if err != nil {
			panic(err)
		}
	}(tmpFile)

	// Copy the BOM-free content to the temporary file
	_, err = io.Copy(tmpFile, bomFreeReader)
	if err != nil {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			return "", err
		}
		return "", err
	}

	return tmpFile.Name(), nil
}
