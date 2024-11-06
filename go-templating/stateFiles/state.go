package stateFiles

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// HandleState is responsible for managing a state file and removing any files that were previously generated but are no longer needed.
// A state file is a specific type of file where each line represents a path, providing a snapshot of the current state of the directory it resides in.
// When new files are generated in the directory, this function updates the state file to reflect those changes.
// This function also provides deletion logic to help remove files that are no longer in the state file.
// Inputs
// - statePath (string): The path to the state file. The state file is a file that contains the desired file paths of the dest directory.
// - exists (bool): Whether or not the statePath (path to the state file) exists.
// - dest (string):  Path to the directory where we handle the state file and deletion of files
// - currStateStr (string): A string that has the file/directory paths in the new state separated by \n
// - ignoreGarbageDeletion (bool): If true, ignores the old state when creating a new state file, so no garbage is removed
// - deletePathList (list): If nil, delete paths that only exist in old state. If not nil, add the paths that only exist in the old state to this list.
// InputsOutputs
// - deletePathList (list): If nil, delete paths that only exist in old state. If not nil, add the paths that only exist in the old state to this list.
// Returns
// - error
func HandleState(statePath string, exists bool, dest string, currStateStr string, ignoreOldState bool, deletePathList *[]string) error {
	// If we have not previously generated the folder, rename temporary file state
	// to be actual hidden file state for current state
	if !exists {
		err := WriteToFile(statePath, currStateStr)
		if err != nil {
			return err
		}
	} else {
		// Create a set for the state
		currState := map[string]bool{}
		currStateSplit := strings.Split(currStateStr, "\n")
		for id := range currStateSplit {
			currState[currStateSplit[id]] = true
		}
		// Retrieve old state
		buf, err := ioutil.ReadFile(statePath)
		if err != nil && !ignoreOldState {
			log.Printf(`Error reading old state file %v. 
Cannot compare old state to new state if old state file does not exist.
Recommended step is to regenerate templates from scratch by deleting generated folder and rerunning generate.
However, if you don't want to regenerate from scratch, simply rerun with --ignoreOldState=true flag, and code will
use new state for state file but garbage may exist from old structure`, err)
			return err
		}
		// Get each file's path from old and new state
		oldState := strings.Split(string(buf), "\n")
		for i := 0; i < len(oldState); i++ {
			// If we have a file in the old state that doesn't exist in the new state
			// remove the file. Do not remove a folder even if it doesn't exist in the new
			// state as it might store user files.
			if currState[oldState[i]] != true {
				destinationPath := dest + "/" + oldState[i]
				info, err := os.Stat(destinationPath)
				if err == nil {
					if !info.IsDir() {
						if deletePathList == nil {
							os.Remove(destinationPath)
						} else {
							*deletePathList = append(*deletePathList, destinationPath)
						}
						// Remove current directory if it is empty
						// TODO: Can optimize as a batch operation at the end rather than test the directory for every deleted file
						currDir := filepath.Dir(destinationPath)
						f, err := os.Open(currDir)
						if err != nil {
							log.Printf("Error reading directory %v \n", err)
							return err
						}
						defer f.Close()
						_, err = f.Readdir(1)
						if err == io.EOF {
							os.RemoveAll(currDir)
						}
					}
				}
			}
		}
		os.Remove(statePath)
		err = WriteToFile(statePath, currStateStr)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteToFile(path string, fileContents string) error {
	// Create file
	_, err := os.Create(path)
	if err != nil {
		log.Printf("Issue creating state file: %v \n", err)
		return err
	}
	// Write current state to file with permissions (0644) such that creator can read and write and
	// others/group can read but not write. No execution permission
	err = os.WriteFile(path, []byte(fileContents), 0644)
	if err != nil {
		log.Printf("Error writing state to file %v \n", err)
		return err
	}
	return nil
}
