package util

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

func CopyFile(directory string, oldFile string, newFile string) {
	r, err := os.Open(oldFile)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	w, err := os.Create(newFile)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		panic(err)
	}
}

func Trim(stream string) string {
	var trimmedString string
	if (strings.Contains(stream, ".cell.balx")) {
		trimmedString = strings.Replace(stream, ".cell.balx", ".celx", -1)
	} else if (strings.Contains(stream, ".bal")) {
		trimmedString = strings.Replace(stream, ".bal", "", -1)
	} else if (strings.Contains(stream, ".cell")) {
		trimmedString = strings.Replace(stream, ".cell", "", -1)
	} else {
		trimmedString = stream
	}
	return trimmedString
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}
