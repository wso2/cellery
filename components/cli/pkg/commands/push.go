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

package commands

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/libtrust"
	"github.com/nokia/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"

	"github.com/cellery-io/sdk/components/cli/pkg/constants"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
)

// RunPush parses the cell image name to recognize the Cellery Registry (A Docker Registry), Organization and version
// and pushes to the Cellery Registry
func RunPush(cellImage string) {
	err := pushImage(cellImage, "", "")
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			username, password, err := util.RequestCredentials()
			if err != nil {
				util.ExitWithErrorMessage("Failed to acquire credentials", err)
			}
			fmt.Println()

			err = pushImage(cellImage, username, password)
			if err != nil {
				util.ExitWithErrorMessage("Failed to push image", err)
			}
		} else {
			util.ExitWithErrorMessage("Failed to pull image", err)
		}
	}
}

func pushImage(cellImage string, username string, password string) error {
	parsedCellImage, err := util.ParseImageTag(cellImage)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while parsing cell image", err)
	}
	repository := parsedCellImage.Organization + "/" + parsedCellImage.ImageName

	spinner := util.StartNewSpinner(fmt.Sprintf("Connecting to %s", util.Bold(parsedCellImage.Registry)))
	defer func() {
		spinner.Stop(true)
	}()

	// Initiating a connection to Cellery Registry
	hub, err := registry.New("https://"+parsedCellImage.Registry, username, password)
	if err != nil {
		spinner.Stop(false)
		util.ExitWithErrorMessage("Error occurred while initializing connection to the Cellery Registry", err)
	}

	imageName := fmt.Sprintf("%s/%s:%s", parsedCellImage.Organization, parsedCellImage.ImageName,
		parsedCellImage.ImageVersion)
	spinner.SetNewAction(fmt.Sprintf("Reading image %s from the Local Repository", util.Bold(imageName)))

	// Reading the cell image
	cellImageFilePath := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, "repo", parsedCellImage.Organization,
		parsedCellImage.ImageName, parsedCellImage.ImageVersion, parsedCellImage.ImageName+constants.CELL_IMAGE_EXT)

	// Checking if the image is present in the local repo
	isImagePresent, _ := util.FileExists(cellImageFilePath)
	if !isImagePresent {
		spinner.Stop(false)
		util.ExitWithErrorMessage(fmt.Sprintf("Failed to push image %s", util.Bold(cellImage)),
			errors.New("Image not Found"))
	}

	cellImageFile, err := os.Open(cellImageFilePath)
	if err != nil {
		spinner.Stop(false)
		util.ExitWithErrorMessage("Error occurred while reading the cell image", err)
	}
	if cellImageFile != nil {
		defer func() {
			err := cellImageFile.Close()
			if err != nil {
				spinner.Stop(false)
				util.ExitWithErrorMessage("Error occurred while opening the cell image", err)
			}
		}()
	}
	cellImageFileBytes, err := ioutil.ReadAll(cellImageFile)
	if err != nil {
		spinner.Stop(false)
		util.ExitWithErrorMessage("Error occurred while reading the cell image", err)
	}

	// Creating the Cell Image Digest (Docker file Layer digest)
	hash := sha256.New()
	hash.Write(cellImageFileBytes)
	sha256sum := hex.EncodeToString(hash.Sum(nil))
	cellImageDigest := digest.NewDigestFromEncoded(digest.SHA256, sha256sum)

	// Checking if the the Cell Image already exists in the registry
	spinner.SetNewAction(fmt.Sprintf("Checking if the image %s already exists in the Registry", util.Bold(imageName)))
	cellImageDigestExists, err := hub.HasBlob(repository, cellImageDigest)
	if err != nil {
		spinner.Stop(false)
		return err
	}

	// Pushing the cell image if it is not already uploaded
	if !cellImageDigestExists {
		spinner.SetNewAction(fmt.Sprintf("Pushing image %s", util.Bold(imageName)))

		// Read stream of files
		err = hub.UploadBlob(repository, cellImageDigest, bytes.NewReader(cellImageFileBytes), nil)
		if err != nil {
			spinner.Stop(false)
			return err
		}
		log.Printf("Successfully uploaded %s cell image", cellImage)
	} else {
		spinner.SetNewAction(fmt.Sprintf("Using already existing image %s in %s Registry", util.Bold(imageName),
			util.Bold(parsedCellImage.Registry)))
		log.Printf("%s cell image already exists", cellImage)
	}

	// Creating a Docker manifest to be uploaded
	cellImageManifest := &schema1.Manifest{
		Name: repository,
		Versioned: manifest.Versioned{
			SchemaVersion: 1,
			MediaType:     schema1.MediaTypeSignedManifest,
		},
		Tag:          parsedCellImage.ImageVersion,
		Architecture: "amd64",
		FSLayers: []schema1.FSLayer{
			{BlobSum: cellImageDigest},
		},
		History: []schema1.History{
			{},
		},
	}

	// Signing the Docker Manifest
	key, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		spinner.Stop(false)
		util.ExitWithErrorMessage("Error occurred while pushing the cell image", err)
	}
	signedCellImageManifest, err := schema1.Sign(cellImageManifest, key)
	if err != nil {
		spinner.Stop(false)
		util.ExitWithErrorMessage("Error occurred while pushing the cell image", err)
	}

	// Uploading the manifest to the Cellery Registry (Docker Registry)
	err = hub.PutManifest(repository, parsedCellImage.ImageVersion, signedCellImageManifest)
	if err != nil {
		spinner.Stop(false)
		return err
	}

	spinner.Stop(true)
	fmt.Print("\n\nImage Digest : " + util.Bold(cellImageDigest))
	util.PrintSuccessMessage(fmt.Sprintf("Successfully pushed cell image: %s", util.Bold(cellImage)))
	util.PrintWhatsNextMessage("pull the image", "cellery pull "+cellImage)

	return nil
}
