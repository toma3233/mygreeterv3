package resourcelinks

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-templating/templateutil"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type deployTmplData struct {
	Name string
	Path string
}
type svcSectionTmplData struct {
	SectionHeader string
	Files         []deployTmplData
}

// GenerateAllResourceFiles traverses through directories to find all desired output files of bicep Dependencies
// For each output file created from a bicep deployment, a corresponding file is generated with the links to the resources
// The function then generates a top-level file that displays all links to these generated files.
// Inputs
// - svcDirName (string): The svcDirName where we want to start the traversal.
// - svcOutputFileName (string): The name of the svc-level file we want to generate.
// Outputs
// - error
func GenerateAllResourceFiles(svcDirName string, svcOutputFileName string) error {
	if _, err := os.Stat(svcDirName); errors.Is(err, os.ErrNotExist) {
		log.Printf("Service path %v does not exist: %v", svcDirName, err)
		return err
	}

	tmplDir := "templates"
	deployTmplFilePath := filepath.Join(tmplDir, "deploy_template.md")
	svcTmplFilePath := filepath.Join(tmplDir, "svc_template.md")

	// Desired file follows the pattern "[.](.*)_output.json".
	rePattern := "[.](.*)_output.json"
	re := regexp.MustCompile(rePattern)

	svcTmplData := make([]svcSectionTmplData, 0)
	deployFileErr := false
	err := filepath.Walk(svcDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error while walking directory: %v", err)
			return err
		}
		if info.IsDir() {
			dirEntries, err := os.ReadDir(path)
			if err != nil {
				log.Printf("Error reading directory %v: %v", path, err)
				return err
			}
			deployTmplDataList := make([]deployTmplData, 0)
			for _, dirEntry := range dirEntries {
				if re.MatchString(dirEntry.Name()) {
					filePath := dirEntry.Name()
					matches := re.FindStringSubmatch(filePath)
					if len(matches) > 1 {
						linkName := matches[1]
						inputFilePath := filepath.Join(path, filePath)
						outputFileName := linkName + "_resources.md"
						outputFilePath := filepath.Join(path, outputFileName)
						linkPath, _ := filepath.Rel(svcDirName, outputFilePath)
						err = createDeploymentFile(inputFilePath, outputFilePath, deployTmplFilePath)
						if err != nil {
							log.Printf("Error creating deployment file: %v", err)
							deployFileErr = true
							continue
						}
						deployTmplDataList = append(deployTmplDataList, deployTmplData{Name: linkName, Path: linkPath})
					}
				}
			}
			if len(deployTmplDataList) > 0 {
				sectionHeader := filepath.Base(path)
				if sectionHeader == filepath.Base(svcDirName) {
					sectionHeader = "Default"
				}
				sectionContent := svcSectionTmplData{
					SectionHeader: sectionHeader,
					Files:         deployTmplDataList,
				}
				svcTmplData = append(svcTmplData, sectionContent)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking directories: %v", err)
		return err
	}
	if (len(svcTmplData) != 1) || (svcTmplData[0].SectionHeader != "Default") {
		outputFilePath := filepath.Join(svcDirName, svcOutputFileName)
		// 0444 grants read access but not write or execute
		err = templateutil.CreateFileFromTemplate(outputFilePath, svcTmplFilePath, svcTmplData, []string{"<<", ">>"}, 0444)
		if err != nil {
			log.Printf("Error creating svc output file %v: %v", outputFilePath, err)
			return err
		}
	}

	if deployFileErr {
		log.Println("\nWARNING: Only some resource-related markdown files successfully created.")
	} else {
		log.Println("\nAll resource-related markdown files successfully created.")
	}
	return nil
}

type tmplInputData struct {
	Resources    []structResourceData
	Dependencies []structDependencyData
}

// For deployment level file
func createDeploymentFile(inputFilePath string, destPath string, tmplFilePath string) error {
	unfmtResourceData, unfmtDependencyData, err := parseJson(inputFilePath)
	if err != nil {
		log.Printf("Error extracting data: %v", err)
		return err
	}

	tmplResourceData := createResourceData(unfmtResourceData)
	tmplDependencyData := createDependencyData(unfmtDependencyData)
	inputData := tmplInputData{
		Resources:    tmplResourceData,
		Dependencies: tmplDependencyData,
	}

	// Allows read access but not write or execute
	err = templateutil.CreateFileFromTemplate(destPath, tmplFilePath, inputData, []string{"<<", ">>"}, 0444)
	if err != nil {
		log.Printf("Error creating svc output file %v: %v", destPath, err)
		return err
	}
	return nil
}

type outputResource struct {
	Id string `json:"id"`
}

type dependencyResource struct {
	Id   string `json:"id"`
	Name string `json:"resourceName"`
}

type dependencyItem struct {
	Id        string               `json:"id"`
	DependsOn []dependencyResource `json:"dependsOn"`
	Name      string               `json:"resourceName"`
}

type properties struct {
	Resources    []outputResource `json:"outputResources"`
	Dependencies []dependencyItem `json:"dependencies"`
}

type deployment struct {
	Properties properties `json:"properties"`
}

func parseJson(filePath string) ([]outputResource, []dependencyItem, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading filepath %v: %v", filePath, err)
		return nil, nil, err
	}
	data := deployment{}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		log.Printf("Error during Unmarshal() of file %v: %v", filePath, err)
		return nil, nil, err
	}
	return data.Properties.Resources, data.Properties.Dependencies, nil
}

type structResourceData struct {
	ResourceName       string
	ResourceType       string
	ResourcePortalLink string
}

func createResourceData(resourceList []outputResource) []structResourceData {
	outputData := make([]structResourceData, len(resourceList))
	for i, resource := range resourceList {
		resourceIdList := strings.Split(resource.Id, "/")
		outputLen := len(resourceIdList)
		name := resourceIdList[outputLen-1]
		portalLink := fmt.Sprintf("https://portal.azure.com/#@microsoft.onmicrosoft.com/resource%v", resource.Id)
		resourceType := strings.Join([]string{resourceIdList[(outputLen - 3)], resourceIdList[(outputLen - 2)]}, "/")
		resourceDataFmt := structResourceData{
			ResourceName:       name,
			ResourceType:       resourceType,
			ResourcePortalLink: portalLink,
		}
		outputData[i] = resourceDataFmt
	}
	return outputData
}

type structDependencyData struct {
	DeploymentName string
	DeploymentId   string
	DependencyList []string
}

func createDependencyData(deploymentList []dependencyItem) []structDependencyData {
	outputData := make([]structDependencyData, len(deploymentList))
	for i, deployment := range deploymentList {
		set := make(map[string]bool)
		var dependencyList []string
		for _, dependency := range deployment.DependsOn {
			if set[dependency.Id] == true {
				continue
			} else {
				set[dependency.Id] = true
				dependencyList = append(dependencyList, dependency.Name)
			}
		}
		newItem := structDependencyData{
			DeploymentName: deployment.Name,
			DeploymentId:   deployment.Id,
			DependencyList: dependencyList,
		}
		outputData[i] = newItem
	}
	return outputData
}
