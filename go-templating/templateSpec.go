package main

import (
	"encoding/csv"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type templateSpecLine struct {
	relativePath string
	overwrite    bool
	users        map[string]bool
	rename       string
}

func generateTemplateSpec(templateFolderStartPath string, options generateOptions) error {
	completeTemplateSpecPath := filepath.Join(templateFolderStartPath, templateSpecPath)
	_, err := os.Stat(completeTemplateSpecPath)
	exists := false
	if err == nil {
		exists = true
	}
	newTemplateSpec := map[string]templateSpecLine{}
	finalSpec := []templateSpecLine{}
	// Walk through all the files/folder in the templates directory.
	err = filepath.Walk(templateFolderStartPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error trying to access certain path %q: %v\n", path, err)
			return err
		}
		// Get relative path based on what starting path is and what our current path is.
		relativePath, err := filepath.Rel(templateFolderStartPath, path)
		// Do not add dashboard templates and template spec to template spec if it exists
		if relativePath != "." && relativePath != templateSpecPath && !strings.Contains(relativePath, nestedTemplatesPath) {
			userMap := map[string]bool{}
			if options.defaultUser {
				userMap = map[string]bool{options.user: true}
			}
			// Set most basic information to allow everything to be overwritten and belonging to aks.
			specLine := templateSpecLine{
				relativePath: relativePath,
				overwrite:    true,
				users:        userMap,
				rename:       "",
			}
			newTemplateSpec[relativePath] = specLine
			finalSpec = append(finalSpec, specLine)
		}
		return nil
	})
	if exists {
		// If template file did previously exist, we will maintain all the
		// settings from previous file for files that still exist in the new spec
		// And add files that were created/ and delete files removed in the new
		// template spec

		// Get old template spec as an array
		oldTemplateSpec, err := readFromSpec(completeTemplateSpecPath)
		if err != nil {
			return err
		}
		// Copy over spec information from old template file for files
		// that still exist.
		for _, val := range oldTemplateSpec {
			_, ok := newTemplateSpec[val.relativePath]
			if ok == true {
				newTemplateSpec[val.relativePath] = val
			}
		}
		// Use final dictionary to adjust sorted spec
		for id, row := range finalSpec {
			finalSpec[id] = newTemplateSpec[row.relativePath]
		}
		// Remove old file and write a new one
		os.Remove(completeTemplateSpecPath)
	}
	err = writeToSpec(completeTemplateSpecPath, finalSpec)
	if err != nil {
		return err
	}
	return nil
}

// Using template spec dictionary write to csv.
func writeToSpec(completeTemplateSpecPath string, templateSpec []templateSpecLine) error {
	// Create template spec file
	templateSpecCSV, err := os.Create(completeTemplateSpecPath)
	if err != nil {
		log.Printf("Failed creating template spec file: %v \n", err)
		return err
	}
	// Write to spec file the contents of the template spec provided
	CSVwriter := csv.NewWriter(templateSpecCSV)
	_ = CSVwriter.Write([]string{"File name", "Overwrite", "Users", "Rename"})
	// It is of utmost importance for the files and folders to be written
	// to the csv in lexicographical as files need to be created in that order
	// for nesting purposes.
	for _, val := range templateSpec {
		var users strings.Builder
		users.WriteString("[")
		var keys []string
		for k := range val.users {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			users.WriteString(key + " ")
		}
		finalUsers := strings.TrimRight(users.String(), " ") + "]"
		// Create a slice of strings for the row
		row := []string{val.relativePath, strconv.FormatBool(val.overwrite), finalUsers, val.rename}
		err = CSVwriter.Write(row)
		if err != nil {
			log.Printf("Can't write to template spec file: %v \n", err)
			return err
		}
	}
	// writing to csv requires flush for data to show up in file
	CSVwriter.Flush()
	templateSpecCSV.Close()
	return nil
}

// Returns an array of information from the csv in sorted format (children come after parent nodes)
func readFromSpec(completeTemplateSpecPath string) ([]templateSpecLine, error) {
	// Retrieve old spec file
	templateSpec := []templateSpecLine{}
	f, err := os.Open(completeTemplateSpecPath)
	if err != nil {
		log.Printf("Template spec cannot be opened: %v \n", err)
		return nil, err
	}

	CSVreader := csv.NewReader(f)
	// Read header Line
	_, err = CSVreader.Read()
	if err != nil {
		log.Printf("Unable to read header from template spec: %v \n", err)
		return nil, err
	}
	// Read all lines except the header in the old template spec
	// and add to dictionary for easy comparison
	for {
		line, err := CSVreader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading from template spec: %v \n", err)
			return nil, err
		}
		overwrite, err := strconv.ParseBool(line[1])
		if err != nil {
			log.Printf("Could not read overwrite from spec %v \n", err)
			return nil, err
		}
		trimmedUsers := strings.TrimLeft(strings.TrimRight(line[2], "]"), "[")
		userList := strings.Split(trimmedUsers, " ")
		userDict := map[string]bool{}
		for _, item := range userList {
			userDict[item] = true
		}
		specLine := templateSpecLine{
			relativePath: line[0],
			overwrite:    overwrite,
			users:        userDict,
			rename:       line[3],
		}
		templateSpec = append(templateSpec, specLine)
	}
	return templateSpec, nil
}
