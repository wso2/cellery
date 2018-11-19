/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

func RecursiveZip(files []string, folders []string, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	for _, folder := range folders {
		err = filepath.Walk(folder, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			relPath := strings.TrimPrefix(filePath, filepath.Dir(folder))

			zipFile, err := myZip.Create(strings.Replace(relPath, "/", "", 1))
			if err != nil {
				return err
			}

			fsFile, err := os.Open(filePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, fsFile)

			if err != nil {
				return err
			}
			return nil
		})
	}

	// Copy files
	for _, file := range files {
		zipFile, err := myZip.Create(file)
		files, err := os.Open(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, files)
	}
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func Unzip(zipFolderName string, destinationFolderName string) error {
	var fileNames []string
	zipFolder, err := zip.OpenReader(zipFolderName)
	if err != nil {
		return err
	}
	defer zipFolder.Close()

	for _, file := range zipFolder.File {
		fileContent, err := file.Open()
		if err != nil {
			return err
		}
		defer fileContent.Close()

		fpath := filepath.Join(destinationFolderName, file.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(destinationFolderName)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		fileNames = append(fileNames, fpath)
		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, fileContent)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func FindInDirectory(directory, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	fileList := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			fileList = append(fileList, filepath.Join(directory, f.Name()))
		}
	}
	return fileList, nil
}
