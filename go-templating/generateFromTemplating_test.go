package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go-templating/stateFiles"
	"path/filepath"

	"io/ioutil"
	"os"
)

var _ = Describe("GenerateFromTemplating", func() {
	var completeStatePath string
	var oldState string
	ignoreGarbageDeletion := false
	dest := "test"
	destStatePath := filepath.Join(dest, statePath)
	// Create ghost test directory that will be considered to be
	// our generated folders
	BeforeEach(func() {
		os.Mkdir(dest, 0777)
		os.Mkdir("test/test1", 0777)
		os.Mkdir("test/test2", 0777)
		os.Mkdir("test/test3", 0777)
		_, _ = os.Create("test/path0.go")
		_, _ = os.Create("test/test1/path1.go")
		_, _ = os.Create("test/test1/path2.go")
		_, _ = os.Create("test/test2/path3.go")
		_, _ = os.Create("test/test3/path4.go")
		_, _ = os.Create("test/test3/path5.go")
		_, _ = os.Create("test/test3/path6.go")
		completeStatePath = filepath.Join("test", statePath)
		oldState = `path0.go
test1
test1/path1.go
test1/path2.go
test2
test2/path3.go
test3
test3/path4.go
test3/path5.go
test3/path6.go`
	})
	AfterEach(func() {
		os.RemoveAll(dest)
	})

	It("Generates state file when it doesn't exist", func() {
		err := stateFiles.HandleState(destStatePath, false, dest, oldState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		state, err := ioutil.ReadFile(completeStatePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(state)).To(Equal(oldState))
	})

	It("State file remains the same when state doesn't change", func() {
		err := stateFiles.HandleState(destStatePath, false, dest, oldState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		newState := oldState
		err = stateFiles.HandleState(destStatePath, true, dest, newState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		state, err := ioutil.ReadFile(completeStatePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(state)).To(Equal(oldState))

	})

	It("Required files are deleted when state no longer has them", func() {
		err := stateFiles.HandleState(destStatePath, false, dest, oldState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		newState := `test1
test1/path1.go
test1/path2.go
test3
test3/path4.go
test3/path6.go`
		err = stateFiles.HandleState(destStatePath, true, dest, newState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		state, err := ioutil.ReadFile(completeStatePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(state)).To(Equal(newState))
		Expect("test/path0.go").ToNot(BeAnExistingFile())
		Expect("test/test1").To(BeADirectory())
		Expect("test/test1/path1.go").To(BeAnExistingFile())
		Expect("test/test1/path2.go").To(BeAnExistingFile())
		Expect("test/test2").ToNot(BeADirectory())
		Expect("test/test2/path3.go").ToNot(BeAnExistingFile())
		Expect("test/test3").To(BeADirectory())
		Expect("test/test3/path4.go").To(BeAnExistingFile())
		Expect("test/test3/path5.go").ToNot(BeAnExistingFile())
		Expect("test/test3/path6.go").To(BeAnExistingFile())
	})

	It("Delete files properly even with extra ones present", func() {
		os.Mkdir("test/test4", 0777)
		_, _ = os.Create("test/test2/test.go")
		_, _ = os.Create("test/test3/path7.go")
		_, _ = os.Create("test/test4/path8.go")
		err := stateFiles.HandleState(destStatePath, false, dest, oldState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		newState := `test1
test1/path1.go
test1/path2.go
test2
test2/test.go
test3
test3/path4.go
test3/path6.go
test3/path7.go
test4
test4/path8.go`
		err = stateFiles.HandleState(destStatePath, true, dest, newState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		state, err := ioutil.ReadFile(completeStatePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(state)).To(Equal(newState))
		Expect("test/path0.go").ToNot(BeAnExistingFile())
		Expect("test/test1").To(BeADirectory())
		Expect("test/test1/path1.go").To(BeAnExistingFile())
		Expect("test/test1/path2.go").To(BeAnExistingFile())
		Expect("test/test2").To(BeADirectory())
		Expect("test/test2/test.go").To(BeAnExistingFile())
		Expect("test/test2/path3.go").ToNot(BeAnExistingFile())
		Expect("test/test3").To(BeADirectory())
		Expect("test/test3/path4.go").To(BeAnExistingFile())
		Expect("test/test3/path5.go").ToNot(BeAnExistingFile())
		Expect("test/test3/path6.go").To(BeAnExistingFile())
		Expect("test/test3/path7.go").To(BeAnExistingFile())
		Expect("test/test4").To(BeADirectory())
		Expect("test/test4/path8.go").To(BeAnExistingFile())

	})

	It("Don't delete extra files added to folder that no longer exists", func() {
		err := stateFiles.HandleState(destStatePath, false, dest, oldState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		// Lets say user creates file under test2
		_, _ = os.Create("test/test2/user.go")
		// Template folders no longer have a test2 folder however
		newState := `path0.go
test1
test1/path1.go
test1/path2.go
test3
test3/path4.go
test3/path5.go
test3/path6.go
`
		err = stateFiles.HandleState(destStatePath, true, dest, newState, ignoreGarbageDeletion, nil)
		Expect(err).NotTo(HaveOccurred())
		state, err := ioutil.ReadFile(completeStatePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(state)).To(Equal(newState))
		Expect("test/path0.go").To(BeAnExistingFile())
		Expect("test/test1").To(BeADirectory())
		Expect("test/test1/path1.go").To(BeAnExistingFile())
		Expect("test/test1/path2.go").To(BeAnExistingFile())
		// We should have deleted file generated by us (path3)
		// But NOT delete the directory if it stores a file
		// that was not created by us.
		Expect("test/test2").To(BeADirectory())
		Expect("test/test2/user.go").To(BeAnExistingFile())
		Expect("test/test2/path3.go").ToNot(BeAnExistingFile())

		Expect("test/test3").To(BeADirectory())
		Expect("test/test3/path4.go").To(BeAnExistingFile())
		Expect("test/test3/path5.go").To(BeAnExistingFile())
		Expect("test/test3/path6.go").To(BeAnExistingFile())
	})
})
