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
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/cellery-io/sdk/components/cli/pkg/constants"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
)

func RunRun(cellImageTag string, instanceName string, withDependencies bool, shareDependencies bool,
	dependencyLinks []string) {
	parsedCellImage, err := util.ParseImageTag(cellImageTag)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while parsing cell image", err)
	}

	imageDir, err := extractImage(parsedCellImage)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while extracting image", err)
	}
	defer func() {
		err = os.RemoveAll(imageDir)
		if err != nil {
			util.ExitWithErrorMessage("Error occurred while cleaning up", err)
		}
	}()

	// Reading Cell Image metadata
	metadataFileContent, err := ioutil.ReadFile(filepath.Join(imageDir, constants.ZIP_ARTIFACTS, "cellery",
		"metadata.json"))
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while reading Cell Image metadata", err)
	}
	cellImageMetadata := &util.CellImageMetaData{}
	err = json.Unmarshal(metadataFileContent, cellImageMetadata)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while reading Cell Image metadata", err)
	}

	if instanceName == "" {
		// Setting a unique instance name since it is not provided
		instanceName, err = generateRandomInstanceName(cellImageMetadata)
		if err != nil {
			util.ExitWithErrorMessage("Error occurred while preparing", err)
		}
	}
	fmt.Printf("\n%s: %s\n", util.Bold("Main Instance"), instanceName)

	// Parsing the dependency links list
	var parsedDependencyLinks []*dependencyAliasLink
	for _, link := range dependencyLinks {
		var dependencyLink *dependencyAliasLink
		linkSplit := strings.Split(link, ":")
		if strings.Contains(linkSplit[0], ".") {
			instanceSplit := strings.Split(linkSplit[0], ".")
			dependencyLink = &dependencyAliasLink{
				Instance:           instanceSplit[0],
				DependencyAlias:    instanceSplit[1],
				DependencyInstance: linkSplit[1],
			}
		} else {
			dependencyLink = &dependencyAliasLink{
				DependencyAlias:    linkSplit[0],
				DependencyInstance: linkSplit[1],
			}
		}
		dependencyLink.IsRunning = isAvailableInRuntime(dependencyLink.DependencyInstance)
		parsedDependencyLinks = append(parsedDependencyLinks, dependencyLink)
	}
	validateDependencyLinks(instanceName, cellImageMetadata, parsedDependencyLinks)

	var mainNode *dependencyTreeNode
	if withDependencies {
		dependencyTree := generateDependencyTree(instanceName, cellImageMetadata, parsedDependencyLinks,
			shareDependencies)
		validateDependencyTree(dependencyTree)
		confirmDependencyTree(dependencyTree)
		startDependencyTree(parsedCellImage.Registry, dependencyTree)
		mainNode = dependencyTree
	} else {
		immediateDependencies := map[string]*dependencyTreeNode{}
		// Check if the provided links are immediate dependencies of the root Cell
		for _, link := range parsedDependencyLinks {
			if !link.IsRunning {
				util.ExitWithErrorMessage("Invalid link",
					fmt.Errorf("all the instances when running without depedencies should be "+
						"running in the runtime, instance %s not available in the runtime", link.DependencyInstance))
			} else if link.Instance == "" || link.Instance == instanceName {
				if _, hasKey := cellImageMetadata.Dependencies[link.DependencyAlias]; hasKey {
					immediateDependencies[link.DependencyAlias] = &dependencyTreeNode{
						Instance:  link.DependencyInstance,
						MetaData:  cellImageMetadata.Dependencies[link.DependencyAlias],
						IsShared:  false,
						IsRunning: link.IsRunning,
					}
				} else {
					var allowedAliases []string
					for alias := range cellImageMetadata.Dependencies {
						allowedAliases = append(allowedAliases, alias)
					}
					util.ExitWithErrorMessage("Invalid links",
						fmt.Errorf("only aliases of the main Cell instance %s: [%s] are allowed when running "+
							"without starting dependencies, received %s", instanceName,
							strings.Join(allowedAliases, ", "), link.DependencyAlias))
				}
			} else {
				util.ExitWithErrorMessage("Invalid links",
					fmt.Errorf("only the main Cell instance %s is allowed when running "+
						"without starting dependencies, received unknown instance %s", instanceName, link.Instance))
			}
		}

		// Check if instances are provided for all the dependencies of the root Cell
		for alias := range cellImageMetadata.Dependencies {
			isLinkProvided := false
			for _, link := range parsedDependencyLinks {
				if link.DependencyAlias == alias {
					isLinkProvided = true
					break
				}
			}
			if !isLinkProvided {
				util.ExitWithErrorMessage("Links for all the dependencies not found",
					fmt.Errorf("required link for alias %s in instance %s not found", alias, instanceName))
			}

		}
		mainNode = &dependencyTreeNode{
			Instance:     instanceName,
			MetaData:     cellImageMetadata,
			IsRunning:    false,
			IsShared:     false,
			Dependencies: immediateDependencies,
		}
		validateDependencyTree(mainNode)
		confirmDependencyTree(mainNode)
	}

	startCellInstance(imageDir, instanceName, mainNode)

	util.PrintSuccessMessage(fmt.Sprintf("Successfully deployed cell image: %s", util.Bold(cellImageTag)))
	util.PrintWhatsNextMessage("list running cells", "cellery list instances")
}

// validateDependencyTree validates the dependency tree of the root instance
func validateDependencyLinks(rootInstance string, rootMetaData *util.CellImageMetaData,
	dependencyLinks []*dependencyAliasLink) {
	// Validating the links provided by the user
	for _, link := range dependencyLinks {
		if link.Instance == "" {
			// This checks for duplicate aliases without parent instance and whether their Cell Images match.
			// If the duplicate aliases have matching Cell Images, then they can share the instance.
			// However, if duplicate aliases are present without parent instances and referring to different
			// Cell Images, the links should be more specific using parent instance
			var validateSubtree func(metadata *util.CellImageMetaData)
			var cellImage *util.CellImage
			validateSubtree = func(metadata *util.CellImageMetaData) {
				for alias, dependencyMetadata := range metadata.Dependencies {
					if alias == link.DependencyAlias {
						if cellImage == nil {
							// This is the first time the alias was found in the dependency tree.
							// (Since the Cell Image was not set)
							cellImage = &util.CellImage{
								Organization: dependencyMetadata.Organization,
								ImageName:    dependencyMetadata.Name,
								ImageVersion: dependencyMetadata.Version,
							}
						} else {
							// Check if the provided alias which is duplicated in the dependency tree is the same image
							if cellImage.Organization != dependencyMetadata.Organization ||
								cellImage.ImageName != dependencyMetadata.Name ||
								cellImage.ImageVersion != dependencyMetadata.Version {
								util.ExitWithErrorMessage("Invalid dependency links",
									fmt.Errorf("duplicated alias %s in dependency tree, referes to different "+
										"images %s/%s:%s and %s/%s:%s, provided aliases which are duplicated in the "+
										"depedencies should be more specific using parent instance", alias,
										cellImage.Organization, cellImage.ImageName, cellImage.ImageVersion,
										dependencyMetadata.Organization, dependencyMetadata.Name,
										dependencyMetadata.Version))
							} else {
								// Since the Cell Image is the same in both aliases the instance will be reused
								fmt.Printf("%s Using a shared instance %s for duplicated alias %s\n",
									util.YellowBold("\U000026A0"), util.Bold(link.DependencyInstance),
									util.Bold(link.DependencyAlias))
							}
						}
					}
					validateSubtree(dependencyMetadata)
				}
			}
			validateSubtree(rootMetaData)
		} else {
			// If the link has a parent instance, this checks if the parent instance had been provided by the user
			// All used parent instances should be specified beforehand as the instance of another alias
			var isLinkParentInstanceProvided bool
			if rootInstance == link.Instance {
				isLinkParentInstanceProvided = true
			} else {
				// Checking if the parent instance in the link is provided as an instance of another alias
				for _, userSpecifiedLink := range dependencyLinks {
					if link.Instance == userSpecifiedLink.DependencyInstance {
						isLinkParentInstanceProvided = true
						break
					}
				}
			}
			if !isLinkParentInstanceProvided {
				util.ExitWithErrorMessage("Error occurred while starting dependencies",
					fmt.Errorf("all parent instances of the provided links should be explicitly given "+
						"as an instance of another alias, instance %s not provided", link.Instance))
			}
		}
	}
}

// generateDependencyOrder reads the metadata and generates a proper start up order for dependencies
func generateDependencyTree(rootInstance string, rootMetaData *util.CellImageMetaData,
	dependencyLinks []*dependencyAliasLink, shareDependencies bool) *dependencyTreeNode {
	// aliasToTreeNodeMap is used to keep track of the already created user provided tree nodes.
	// The key of the is the alias and the value is the tree node.
	// The alias is prefixed by the instance only if the user specified the parent instance as well.
	aliasToTreeNodeMap := map[string]*dependencyTreeNode{}

	// generatedInstanceTreeNodes are used to keep track of the instances automatically generated
	// These will be shared among the auto generated instances based on "shareDependencies" environment variable
	var generatedInstanceTreeNodes []*dependencyTreeNode

	// traverseDependencies traverses through the dependency tree and populates the startup order considering the
	// relationship between dependencies
	var traverseDependencies func(instance string, metaData *util.CellImageMetaData, treeNode *dependencyTreeNode)
	traverseDependencies = func(instance string, metaData *util.CellImageMetaData, treeNode *dependencyTreeNode) {
		for alias, dependencyMetaData := range metaData.Dependencies {
			var dependencyNode *dependencyTreeNode

			// Check if the dependency link is provided by the user
			for _, link := range dependencyLinks {
				if alias == link.DependencyAlias && (link.Instance == "" || link.Instance == instance) {
					var aliasPrefix string
					if link.Instance != "" {
						aliasPrefix = link.Instance + "."
					}
					key := aliasPrefix + alias

					if node, hasKey := aliasToTreeNodeMap[key]; hasKey {
						dependencyNode = node
						dependencyNode.IsShared = true
					} else {
						dependencyNode = &dependencyTreeNode{
							Instance:     link.DependencyInstance,
							MetaData:     dependencyMetaData,
							Dependencies: map[string]*dependencyTreeNode{},
							IsShared:     false,
							IsRunning:    link.IsRunning,
						}
						aliasToTreeNodeMap[key] = dependencyNode
					}
					break
				}
			}

			if dependencyNode == nil {
				if shareDependencies {
					// Check if an instance had been already allocated for this image
					for _, existingNode := range generatedInstanceTreeNodes {
						if existingNode.MetaData.Organization == dependencyMetaData.Organization &&
							existingNode.MetaData.Name == dependencyMetaData.Name &&
							existingNode.MetaData.Version == dependencyMetaData.Version {
							dependencyNode = existingNode
							existingNode.IsShared = true
						}
					}
				}

				if dependencyNode == nil {
					dependencyInstance, err := generateRandomInstanceName(dependencyMetaData)
					if err != nil {
						util.ExitWithErrorMessage("Error occurred while starting dependencies", err)
					}
					dependencyNode = &dependencyTreeNode{
						Instance:     dependencyInstance,
						MetaData:     dependencyMetaData,
						Dependencies: map[string]*dependencyTreeNode{},
						IsShared:     false,
						IsRunning:    false,
					}
					generatedInstanceTreeNodes = append(generatedInstanceTreeNodes, dependencyNode)
				}
			}

			// Traversing the dependencies (Depth First Search for start up items)
			treeNode.Dependencies[alias] = dependencyNode
			traverseDependencies(dependencyNode.Instance, dependencyMetaData, dependencyNode)
		}
	}
	dependencyTreeRoot := &dependencyTreeNode{
		Instance:     rootInstance,
		MetaData:     rootMetaData,
		Dependencies: map[string]*dependencyTreeNode{},
		IsShared:     false,
		IsRunning:    false,
	}
	traverseDependencies(rootInstance, rootMetaData, dependencyTreeRoot)
	return dependencyTreeRoot
}

// validateDependencyTree validates a generated dependency tree
func validateDependencyTree(treeRoot *dependencyTreeNode) {
	// Validate whether the Cell Image of all the specified instances match
	instanceToMetadataMap := map[string]*util.CellImageMetaData{}
	var validateDependencySubtreeOffline func(subTreeRoot *dependencyTreeNode)
	validateDependencySubtreeOffline = func(subTreeRoot *dependencyTreeNode) {
		for _, dependency := range subTreeRoot.Dependencies {
			if metadata, hasKey := instanceToMetadataMap[dependency.Instance]; hasKey {
				if metadata.Organization != dependency.MetaData.Organization ||
					metadata.Name != dependency.MetaData.Name ||
					metadata.Version != dependency.MetaData.Version {
					util.ExitWithErrorMessage("Invalid instance linking",
						fmt.Errorf("instance %s cannot be shared by different Cell Images %s/%s:%s and %s/%s:%s",
							dependency.Instance,
							dependency.MetaData.Organization, dependency.MetaData.Name, dependency.MetaData.Version,
							metadata.Organization, metadata.Name, metadata.Version))
				}
			} else {
				instanceToMetadataMap[dependency.Instance] = dependency.MetaData
			}
			validateDependencySubtreeOffline(dependency)
		}
	}
	instanceToMetadataMap[treeRoot.Instance] = treeRoot.MetaData
	validateDependencySubtreeOffline(treeRoot)

	// Validating whether the instances running in the runtime match the linked image
	// TODO : Add logic to call the runtime
}

// confirmDependencyTree confirms from the user whether the intended dependency tree had been resolved
func confirmDependencyTree(tree *dependencyTreeNode) {
	var dependencyData [][]string
	var traversedInstances []string
	var extractDependencyTreeData func(subTree *dependencyTreeNode)
	extractDependencyTreeData = func(subTree *dependencyTreeNode) {
		for _, dependency := range subTree.Dependencies {
			// Traversing the dependency tree
			if !dependency.IsRunning {
				extractDependencyTreeData(dependency)
			}

			// Adding used instances table content
			instanceAlreadyAdded := false
			for _, instance := range traversedInstances {
				if instance == dependency.Instance {
					instanceAlreadyAdded = true
					break
				}
			}
			if !instanceAlreadyAdded {
				var usedInstance string
				if dependency.IsRunning {
					usedInstance = "Available in Runtime"
				} else {
					usedInstance = "To be Created"
				}
				var sharedSymbol string
				if dependency.IsShared {
					sharedSymbol = "Shared"
				} else {
					sharedSymbol = " - "
				}
				dependencyData = append(dependencyData, []string{
					dependency.Instance,
					dependency.MetaData.Organization + "/" + dependency.MetaData.Name + ":" + dependency.MetaData.Version,
					usedInstance,
					sharedSymbol,
				})
				traversedInstances = append(traversedInstances, dependency.Instance)
			}
		}
	}
	extractDependencyTreeData(tree)
	dependencyData = append(dependencyData, []string{
		tree.Instance,
		tree.MetaData.Organization + "/" + tree.MetaData.Name + ":" + tree.MetaData.Version,
		"To be Created",
		" - ",
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"INSTANCE NAME", "CELL IMAGE", "USED INSTANCE", "SHARED"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAlignment(3)
	table.SetRowSeparator("-")
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{})
	table.AppendBulk(dependencyData)
	fmt.Printf("\n%s:\n\n", util.Bold("Instances to be Used"))
	table.Render()

	fmt.Printf("\n%s:\n\n %s\n", util.Bold("Dependency Tree to be Used"), util.Bold(tree.Instance))
	var printDependencyTree func(subTree *dependencyTreeNode, nestingLevel int, ancestorBranchPrintRequirement []bool)
	printDependencyTree = func(subTree *dependencyTreeNode, nestingLevel int, ancestorBranchPrintRequirement []bool) {
		var index = 0
		for alias, dependency := range subTree.Dependencies {
			// Adding the dependency tree visualization content
			for j := 0; j < nestingLevel; j++ {
				if ancestorBranchPrintRequirement[j] {
					fmt.Print("   │ ")
				} else {
					fmt.Print("     ")
				}
			}
			if index == len(subTree.Dependencies)-1 {
				fmt.Print("   └")
			} else {
				fmt.Print("   ├")
			}
			fmt.Printf("── %s: %s\n", util.Bold(alias), dependency.Instance)

			// Traversing the dependency tree
			if !dependency.IsRunning {
				printDependencyTree(dependency, nestingLevel+1,
					append(ancestorBranchPrintRequirement, index != len(subTree.Dependencies)-1))
			}
			index++
		}
	}
	printDependencyTree(tree, 0, []bool{})
	fmt.Println()

	fmt.Printf("%s Do you wish to continue with starting above Cell instances (Y/n)? ", util.YellowBold("?"))
	reader := bufio.NewReader(os.Stdin)
	confirmation, err := reader.ReadString('\n')
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while confirming the dependency tree", err)
	}
	if strings.ToLower(strings.TrimSpace(confirmation)) == "n" {
		fmt.Printf("%s Running Cell aborted\n", util.Red("\U0000274C"))
		os.Exit(0)
	}
	fmt.Println()
}

// startDependencyTree starts up the whole dependency tree except the root
// This does not start the root of the dependency tree
func startDependencyTree(registry string, tree *dependencyTreeNode) {
	var wg sync.WaitGroup
	wg.Add(len(tree.Dependencies))
	for _, dependency := range tree.Dependencies {
		if dependency.IsRunning {
			wg.Done()
		} else { // This level of checking is done to not start unwanted goroutines
			go func(dependencyNode *dependencyTreeNode) {
				defer wg.Done()
				dependencyNode.Mux.Lock()
				if !dependencyNode.IsRunning { // This level of checking is done to make sure the condition is met
					startDependencyTree(registry, dependencyNode)
					//cellImage := &util.CellImage{
					//	Registry:     registry,
					//	Organization: dependencyNode.MetaData.Organization,
					//	ImageName:    dependencyNode.MetaData.Name,
					//	ImageVersion: dependencyNode.MetaData.Version,
					//}
					//imageDir, err := extractImage(cellImage)
					//if err != nil {
					//	util.ExitWithErrorMessage("Error occurred while extracting image", err)
					//}

					startCellInstance("", dependencyNode.Instance, dependencyNode)
					dependencyNode.IsRunning = true

					//err = os.RemoveAll(imageDir)
					//if err != nil {
					//	util.ExitWithErrorMessage("Error occurred while cleaning up", err)
					//}
				}
				dependencyNode.Mux.Unlock()
			}(dependency)
		}
	}
	wg.Wait()
}

// extractImage extracts the image into a temporary directory and returns the path.
// Cleaning the path after finishing your work is your responsibility.
func extractImage(cellImage *util.CellImage) (string, error) {
	repoLocation := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, "repo", cellImage.Organization,
		cellImage.ImageName, cellImage.ImageVersion)
	zipLocation := filepath.Join(repoLocation, cellImage.ImageName+constants.CELL_IMAGE_EXT)
	cellImageTag := cellImage.Organization + "/" + cellImage.ImageName + ":" + cellImage.ImageVersion

	// Pull image if not exist
	imageExists, err := util.FileExists(zipLocation)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while extracting the Cell image", err)
	}
	if !imageExists {
		fmt.Printf("\nUnable to find image %s locally.", cellImageTag)
		fmt.Printf("\nPulling image: %s", cellImageTag)
		RunPull(cellImageTag, true)
	}

	// Unzipping image to a temporary location
	tempPath, err := ioutil.TempDir("", "cellery-cell-image")
	if err != nil {
		return "", err
	}
	err = util.Unzip(zipLocation, tempPath)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while extracting cell image", err)
	}
	return tempPath, nil
}

func isAvailableInRuntime(instance string) bool {
	// TODO : Check if the instance is available in the runtime
	return instance == "test-inst" || instance == "stock-inst" || instance == "employee-inst"
}

func startCellInstance(imageDir string, instanceName string, runningNode *dependencyTreeNode) {
	// TODO : Remove this method and rename the below method to match this method
	fmt.Println("Started " + instanceName)
	time.Sleep(2 * time.Second)
	fmt.Println("Running " + instanceName)
}

func startCellInstance1(imageDir string, instanceName string, runningNode *dependencyTreeNode) {
	balFileName, err := util.GetSourceFileName(filepath.Join(imageDir, constants.ZIP_BALLERINA_SOURCE))
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while extracting source file: ", err)
	}
	balFilePath := filepath.Join(imageDir, constants.ZIP_BALLERINA_SOURCE, balFileName)

	containsRunFunction, err := util.RunMethodExists(balFilePath)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while checking for run function ", err)
	}
	if containsRunFunction {
		// Generating the first level dependency map
		dependencies := map[string]string{}
		for alias, dependency := range runningNode.Dependencies {
			dependencies[alias] = dependency.Instance
		}

		// Calling the run function
		dependenciesJson, err := json.Marshal(dependencies)
		if err != nil {
			util.ExitWithErrorMessage("Error occurred while starting Cell instance"+instanceName, err)
		}
		cmd := exec.Command("ballerina", "run", balFilePath+":run",
			runningNode.MetaData.Organization+"/"+runningNode.MetaData.Name, runningNode.MetaData.Version,
			instanceName, string(dependenciesJson))
		cmd.Env = append(cmd.Env, constants.CELLERY_IMAGE_DIR_ENV_VAR+"="+imageDir)
		stdoutReader, _ := cmd.StdoutPipe()
		stdoutScanner := bufio.NewScanner(stdoutReader)
		go func() {
			for stdoutScanner.Scan() {
				fmt.Printf("\033[36m%s\033[m\n", stdoutScanner.Text())
			}
		}()
		stderrReader, _ := cmd.StderrPipe()
		stderrScanner := bufio.NewScanner(stderrReader)
		go func() {
			for stderrScanner.Scan() {
				fmt.Printf("\033[36m%s\033[m\n", stderrScanner.Text())
			}
		}()
		err = cmd.Start()
		if err != nil {
			util.ExitWithErrorMessage("Error in executing cellery run", err)
		}
		err = cmd.Wait()
		if err != nil {
			util.ExitWithErrorMessage("Error occurred in cellery run", err)
		}
	}

	// Update the Cell instance name
	celleryDir := filepath.Join(imageDir, constants.ZIP_ARTIFACTS, "cellery")
	k8sYamlFile := filepath.Join(celleryDir, runningNode.MetaData.Name+".yaml")
	if instanceName != "" {
		// Cell instance name changed.
		err = util.ReplaceInFile(k8sYamlFile, "name: "+runningNode.MetaData.Name, "name: "+instanceName, 1)
		if err != nil {
			util.ExitWithErrorMessage("Error in replacing cell instance name", err)
		}
	}

	// Applying the yaml
	cmd := exec.Command(constants.KUBECTL, constants.APPLY, "-f", celleryDir)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Printf("\033[36m%s\033[m\n", stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Printf("\033[36m%s\033[m\n", stderrScanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		util.ExitWithErrorMessage("Error in executing cell run", err)
	}
	err = cmd.Wait()
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while running cell image", err)
	}
}

// newUUID generates a random UUID according to RFC 4122
func generateRandomInstanceName(dependencyMetaData *util.CellImageMetaData) (string, error) {
	u := make([]byte, 4)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x", u)

	// Generating random instance name
	return dependencyMetaData.Organization + "-" + dependencyMetaData.Name + "-" +
		strings.Replace(dependencyMetaData.Version, ".", "-", -1) + "-" + uuid, nil
}

// dependencyAliasLink is used to store the link information provided by the user
type dependencyAliasLink struct {
	Instance           string
	DependencyAlias    string
	DependencyInstance string
	IsRunning          bool
}

// dependencyTreeNode is used as a node of the dependency tree
type dependencyTreeNode struct {
	Mux          sync.Mutex
	Instance     string
	MetaData     *util.CellImageMetaData
	Dependencies map[string]*dependencyTreeNode
	IsShared     bool
	IsRunning    bool
}
