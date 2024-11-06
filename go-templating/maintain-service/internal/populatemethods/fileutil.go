package populatemethods

import (
	"errors"
	"go-templating/stateFiles"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go-templating/templateutil"

	"github.com/emicklei/proto"
)

type methodDirArgs struct {
	destinationDirPath string
	templateName       string
	outputFileType     string
}

func PopulateAllMethodFiles(svcDirName string, protoFilePath string, deleteGarbageFiles bool) error {
	if _, err := os.Stat(protoFilePath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Proto file path %s does not exist: %v", protoFilePath, err)
		return err
	}

	methodsList, err := parseMethods(protoFilePath)
	if err != nil {
		log.Printf("Error creating method list: %v", err)
		return err
	}

	var deletePathList []string

	// Method templates follow the pattern ".method_template_<filetype>.txt".
	var templateName string
	rePattern := ".method_template_([a-zA-Z]+).txt"
	re := regexp.MustCompile(rePattern)

	root := svcDirName
	if _, err := os.Stat(svcDirName); errors.Is(err, os.ErrNotExist) {
		log.Printf("Service path %s does not exist: %v", svcDirName, err)
		return err
	}

	// filepath.Walk does not automatically update for deleted files during traversal.
	// Deleting files during filepath.Walk cause an error. Adding files during filepath.Walk is okay.
	tmplFileErr := false
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error while walking directory: %v", err)
			return err
		}
		if info.IsDir() {
			dirEntries, _ := os.ReadDir(path)
			for _, dirEntry := range dirEntries {
				if re.MatchString(dirEntry.Name()) {
					templateName = dirEntry.Name()
					matches := re.FindStringSubmatch(dirEntry.Name()) // Parsing the file name for output file type.
					if len(matches) > 0 {
						outputFileType := "." + matches[1]
						dirArgs := methodDirArgs{
							destinationDirPath: path,
							templateName:       templateName,
							outputFileType:     outputFileType,
						}
						err = populateMethodFilesInTemplateDir(dirArgs, methodsList, &deletePathList, svcDirName)
						if err != nil {
							log.Printf("Error using template %s and populating files in path %s: %v", templateName, path, err)
							tmplFileErr = true
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking directories: %v", err)
		return err
	}
	if deleteGarbageFiles {
		for _, path := range deletePathList {
			err = os.Remove(path)
			if err != nil {
				log.Printf("Error deleting file %s: %v", path, err)
				return err
			}
			log.Printf("%s deleted.", path)
		}
	}
	if tmplFileErr {
		log.Println("\nWARNING: Only some proto method-related successfully created.")

	} else {
		log.Println("\nAll proto method-related files successfully created.")
	}
	return nil
}

// Parses the methods from the proto file and creates a list of the methods
func parseMethods(protoPath string) ([]*proto.RPC, error) {
	var methodsList []*proto.RPC
	reader, err := os.Open(protoPath)

	if err != nil {
		log.Printf("Error opening proto file: %v", err)
		return nil, err
	}

	defer reader.Close()

	parser := proto.NewParser(reader)

	definition, err := parser.Parse()

	if err != nil {
		log.Printf("Error parsing file: %v", err)
		return nil, err
	}

	proto.Walk(definition,
		proto.WithRPC(func(r *proto.RPC) {
			methodsList = append(methodsList, r)
		}),
	)
	return methodsList, nil
}

func populateMethodFilesInTemplateDir(m methodDirArgs, methodsList []*proto.RPC, deletePathList *[]string, svcDirName string) error {
	templatePath := filepath.Join(m.destinationDirPath, m.templateName)
	var currStateStr strings.Builder // Used to keep track of current method files.
	numMethods := len(methodsList)
	for i := 0; i < numMethods; i++ {
		var newFileName string
		// If the output file is a bicep file, the file name should be appended with ".ServiceResources.Template" to match Ev2 naming conventions.
		if m.outputFileType == ".bicep" {
			newFileName = strings.Join([]string{methodsList[i].Name, ".ServiceResources.Template", m.outputFileType}, "")
		} else {
			newFileName = strings.Join([]string{methodsList[i].Name, m.outputFileType}, "")
		}
		destinationPath := filepath.Join(m.destinationDirPath, newFileName)
		relativePathForFile := newFileName + "\n"
		currStateStr.WriteString(relativePathForFile)
		_, err := os.Stat(destinationPath)
		relPath, _ := filepath.Rel(svcDirName, destinationPath)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Error checking file %s: %v\n", relPath, err)
			return err
		} else if err == nil {
			log.Printf("%s exists, no changes made.", relPath)
			continue
		} else if os.IsNotExist(err) {
			// 0666 grants full read and write access to anyone
			err = templateutil.CreateFileFromTemplate(destinationPath, templatePath, methodsList[i], []string{"{{", "}}"}, 0666)
			if err != nil {
				log.Printf("Error creating file from template %v: %v", destinationPath, err)
				return err
			}
			log.Printf("%s created.", relPath)
		}
	}
	statePath := filepath.Join(m.destinationDirPath, "/.methods_state.txt")
	if m.destinationDirPath == "" {
		statePath = ".methods_state.txt"
	}
	exists := pathExists(statePath)
	stateFiles.HandleState(statePath, exists, m.destinationDirPath, currStateStr.String(), false, deletePathList)
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
