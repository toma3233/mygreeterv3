package templateutil

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"gopkg.in/yaml.v3"
)

// ExecuteTemplatesInFolder traverses through directories to find all desired template files that we want to execute
// data into. The files are idenitfied using the "templatePrefix" input, and the data that is used to execute the template
// is found at the path "envConfig"
// Inputs
// - templatePrefix (string): The prefix used to identify template files.
// - folderPath (string): The path for the folder we are traversing.
// - envConfig (string): The path to the data yaml file.
// - delimiter ([]string): The delimiters we are using for template vars.
// Outputs
// - error
func ExecuteTemplatesInFolder(templatePrefix string, folderPath string, envConfig string, delimiter []string) error {
	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Folder Path %v does not exist: %v", folderPath, err)
		return err
	}
	data := map[string]string{}
	conf, err := ioutil.ReadFile(envConfig)
	if err != nil {
		log.Printf("Could not read env config file: %v \n", err)
		return err
	}
	err = yaml.Unmarshal(conf, &data)
	if err != nil {
		log.Printf("Could not unmarshal env config file: %v \n", err)
		return err
	}
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error while walking directory: %v", err)
			return err
		}
		//Only consider files and not directories
		if !info.IsDir() {
			fileName := filepath.Base(path)
			if strings.HasPrefix(fileName, templatePrefix) {
				destinationPath := filepath.Join(filepath.Dir(path), fileName[len(templatePrefix):])
				// 0666 grants full read and write access to anyone
				err = CreateFileFromTemplate(destinationPath, path, data, delimiter, 0666)
				if err != nil {
					log.Printf("Could not create the file from template: %v \n", err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Could not walk through folder: %v \n", err)
		return err
	}
	return nil
}

// Creates a file from the inputted template file.
// Inputs
// - outputFilePath (string): The file path of the generated file.
// - tmplFilePath (string): The path to the template file.
// - tmplData (any): The data used to populate the template file.
// - delims (string array, length 2): The delimiters used in the template file (i.e. []string{"<<", ">>"}).
// - filePermissions (fs.FileMode): The file permissions of the generated file (i.e. 0777).
// Outputs
// - err
func CreateFileFromTemplate(outputFilePath string, tmplFilePath string, tmplData any, delims []string, filePermissions fs.FileMode) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Error creating file %v: %v", outputFilePath, err)
		return err
	}
	defer func() error {
		if err := outputFile.Close(); err != nil {
			log.Printf("Error closing file %v: %v \n", outputFilePath, err)
			return err
		}
		return nil
	}()
	tmplName := filepath.Base(tmplFilePath)
	if len(delims) != 2 {
		log.Printf("Incorrect number of delims passed in.")
	}
	funcMap := CreateFuncMap()
	resourceTmpl, err := template.New(tmplName).Delims(delims[0], delims[1]).Option("missingkey=error").Funcs(funcMap).ParseFiles(tmplFilePath)
	if err != nil {
		log.Printf("Issue parsing template file at %v, %v \n", tmplFilePath, err)
		os.Remove(outputFilePath)
		return err
	}
	err = resourceTmpl.ExecuteTemplate(outputFile, tmplName, tmplData)
	if err != nil {
		log.Printf("Template execution error: %v, %v \n", outputFilePath, err)
		os.Remove(outputFilePath)
		return err
	}
	err = os.Chmod(outputFilePath, filePermissions)
	if err != nil {
		log.Printf("Error setting permissions: %v \n", err)
		return err
	}
	return nil
}

// Creates a function map used when templating to perform functions on data variables.
// Outputs
// - template.FuncMap
func CreateFuncMap() template.FuncMap{
	funcMap := template.FuncMap{
		// Contains is used to check if goModuleNamePrefix uses go.goms.io or dev.azure.com
		"contains": strings.Contains,
		"upper":    strings.ToUpper,
		// Remove git is required as for git configuration in user environment, goModuleNamePrefix is used
		// without the git suffix.
		"trimGitSuffix": func(input string) string {
			index := strings.Index(input, ".git")
			if index != -1 {
				return input[:index]
			}
			return input
		},
		"apiModule": func(goModuleNamePrefix string, serviceDirectoryName string) string {
			return filepath.Join(filepath.Join(goModuleNamePrefix, serviceDirectoryName), "api")
		},
		"serverModule": func(goModuleNamePrefix string, serviceDirectoryName string) string {
			return filepath.Join(filepath.Join(goModuleNamePrefix, serviceDirectoryName), "server")
		},
	}
	return funcMap
}