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
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CopyFile(oldFile string, newFile string) {
	input, err := os.Open(oldFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	output, err := os.Create(newFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
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

			zipFile, err := myZip.Create(relPath)
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

func FindInDirectory(directory, suffix string) ([]string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil
	}
	fileList := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			fileList = append(fileList, filepath.Join(directory, f.Name()))
		}
	}
	return fileList
}

func GetDuration(startTime time.Time) string {
	duration := ""
	var year, month, day, hour, min, sec int
	currentTime := time.Now()
	if startTime.Location() != currentTime.Location() {
		currentTime = currentTime.In(startTime.Location())
	}
	if startTime.After(currentTime) {
		startTime, currentTime = currentTime, startTime
	}
	startYear, startMonth, startDay := startTime.Date()
	currentYear, currentMonth, currentDay := currentTime.Date()

	startHour, startMinute, startSecond := startTime.Clock()
	currentHour, currentMinute, currentSecond := currentTime.Clock()

	year = int(currentYear - startYear)
	month = int(currentMonth - startMonth)
	day = int(currentDay - startDay)
	hour = int(currentHour - startHour)
	min = int(currentMinute - startMinute)
	sec = int(currentSecond - startSecond)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(startYear, startMonth, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	numOfTimeUnits := 0
	if year > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(year) + " years "
		numOfTimeUnits++
	}
	if month > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(month) + " months "
		numOfTimeUnits++
	}
	if day > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(day) + " days "
		numOfTimeUnits++
	}
	if hour > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(hour) + " hours "
		numOfTimeUnits++
	}
	if min > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(min) + " minutes "
		numOfTimeUnits++
	}
	if sec > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(sec) + " seconds"
		numOfTimeUnits++
	}
	return duration
}

func ConvertStringToTime(timeString string) time.Time {
	convertedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error parsing time: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	return convertedTime
}

// Creates a new file upload http request with optional extra params
func FileUploadRequest(uri string, params map[string]string, paramName, path string, secure bool) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	if (!secure) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
