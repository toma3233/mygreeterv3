package main

import (
	"go-templating/stateFiles"
	"go-templating/templateutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func genFromTemplate(allInput allInput, options generateOptions) error {
	// Destination directory and templates path must be changed
	// if generating middleware instead of service
	dest, templateName := getTemplateVars(allInput, options.generationType)
	start := allInput.GeneratorInputs[templateName]
	allInput.User = options.user
	allInputMap := ConvertStructToMap(allInput)
	// We do not want to pass the user to the template files information
	allInputMap["generatorInputs"] = map[string]string{}
	// Create destination directory if it doesn't exist
	err := os.Mkdir(dest, 0777)
	exists := os.IsExist(err)
	if err != nil && !os.IsExist(err) {
		log.Printf("Could not create directory: %v \n", err)
		return err
	}
	// Create a temporary directory within destination to store files.
	tempFolder := filepath.Join(os.TempDir(), "tempDir")
	err = os.Mkdir(tempFolder, 0777)
	if err != nil && !os.IsExist(err) {
		log.Printf("Could not create temporary directory: %v \n", err)
		return err
	}
	// Create a func map that provides options to perform functions on data variables.
	funcMap := templateutil.CreateFuncMap()
	// Create current state as a string
	var currStateStr strings.Builder
	completeTemplateSpecPath := filepath.Join(start, templateSpecPath)
	templateSpec, err := readFromSpec(completeTemplateSpecPath)
	if err != nil {
		return err
	}
	var nestedTemplate *template.Template
	var nested bool
	for _, specInfo := range templateSpec {
		if specInfo.users[allInput.User] == true {
			templatePath := filepath.Join(start, specInfo.relativePath)
			info, err := os.Stat(templatePath)
			if err != nil {
				log.Printf("The template spec path does not exist in the templates folder : %v, %v \n", templatePath, err)
				os.RemoveAll(tempFolder)
				return err
			}
			temporaryPath := filepath.Join(tempFolder, specInfo.relativePath)
			destinationPath := filepath.Join(dest, specInfo.relativePath)
			destinationPathForState := specInfo.relativePath + "\n"
			if specInfo.rename != "" {
				temporaryPath = filepath.Join(tempFolder, specInfo.rename)
				destinationPath = filepath.Join(dest, specInfo.rename)
				destinationPathForState = specInfo.rename + "\n"
			}
			// Add relative path to current state
			currStateStr.WriteString(destinationPathForState)
			// If it is a directory within the templates folder, create directory in destination directory
			if info.IsDir() {
				if specInfo.relativePath != "." {
					err = os.Mkdir(temporaryPath, 0777)
					// Do not overwrite folder if it already exists
					if err != nil && !os.IsExist(err) {
						// TODO: If a file is deleted and then the same name is used as a directory in the template,
						// the directory creation will fail because the same name is already used.
						log.Printf("Issue creating nested directory: %v \n", err)
						os.RemoveAll(tempFolder)
						return err
					}
				}
				nestedPath := filepath.Join(templatePath, nestedTemplatesPath)
				_, err = os.Stat(nestedPath)
				// If nested template folder exists, add files to
				// a list so they can be parsed later
				if err == nil {
					nested = true
					items, _ := os.ReadDir(nestedPath)
					var fileList []string
					for _, item := range items {
						path := filepath.Join(nestedPath, item.Name())
						fileList = append(fileList, path)
					}
					nestedTemplate, err = template.New("").Delims("<<", ">>").ParseFiles(fileList...)
					if err != nil {
						log.Printf("Issue parsing templates files in directory %v, %v \n", fileList, err)
						os.RemoveAll(tempFolder)
						return err
					}
				} else {
					nested = false
				}
			} else {
				// If it is a file, check whether or not it exist
				_, err = os.Stat(destinationPath)
				if err != nil && !(os.IsNotExist(err)) {
					log.Printf("Error with destination path file : %v, %v \n", destinationPath, err)
					os.RemoveAll(tempFolder)
					return err
				}
				if (specInfo.overwrite) || (os.IsNotExist(err)) {
					// If the file already exists and we have overwrite capabilities,
					// Parse it, create a new one, and execute the template into the new file.
					// Otherwise do nothing
					var tmp *template.Template
					currFile, err := os.Create(temporaryPath)
					if err != nil {
						log.Printf("Issue creating file: %v \n", err)
						os.RemoveAll(tempFolder)
						return err
					}
					fileName := filepath.Base(templatePath)
					if nested {
						tmp, err = nestedTemplate.Clone()
						if err != nil {
							log.Printf("Issue cloning nested templates  %v \n", err)
							os.RemoveAll(tempFolder)
							return err
						}
						// Parse all the nested templates
						tmp, err = tmp.New(fileName).Delims("<<", ">>").Option("missingkey=error").Funcs(funcMap).ParseFiles(templatePath)
					} else {
						tmp, err = template.New(fileName).Delims("<<", ">>").Option("missingkey=error").Funcs(funcMap).ParseFiles(templatePath)
					}
					if err != nil {
						log.Printf("Issue parsing template file at %v, %v \n", templatePath, err)
						os.RemoveAll(tempFolder)
						return err
					}
					// Make sure scripts are executable.
					if strings.Contains(temporaryPath, ".sh") {
						os.Chmod(temporaryPath, 0755)
					}
					// Execute templates into new file in destination directory with provided variables
					// for service vs middleware.
					err = tmp.ExecuteTemplate(currFile, fileName, allInputMap)
					if err != nil {
						log.Printf("Template execution error: %v, %v \n", temporaryPath, err)
						os.RemoveAll(tempFolder)
						return err
					}
					// TODO: Reformat file specific to file type
					currFile.Close()
				}
			}
		}
	}
	// Copy over temporary folder to destination and remove temporary folder.
	allTemp := tempFolder + string(os.PathSeparator) + "."
	cmd := exec.Command("cp", "-a", allTemp, dest)
	_, err = cmd.Output()
	os.RemoveAll(tempFolder)
	if err != nil {
		log.Printf("Error copying temp folder to destination: %v \n", err)
		return err
	}
	completeStatePath := filepath.Join(dest, statePath)
	if dest == "" {
		completeStatePath = statePath
	}
	err = stateFiles.HandleState(completeStatePath, exists, dest, currStateStr.String(), options.ignoreOldState, nil)
	if err != nil {
		return err
	}
	return nil
}
