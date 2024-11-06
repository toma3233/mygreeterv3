package main

import (
	"go-templating/stateFiles"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateSpec testing", func() {
	templateFolderStartPath := "test"
	var completeTemplateSpecPath string
	var oldTemplateSpec string
	destinationPath := "dest"
	var testOptions generateOptions
	var allInput allInput
	// Create ghost test directory that will be considered to be
	// our template folders
	BeforeEach(func() {
		os.Mkdir(templateFolderStartPath, 0777)
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
		testOptions = generateOptions{}
		testOptions.defaultUser = true
		testOptions.user = "external"
		testOptions.ignoreOldState = false
		completeTemplateSpecPath = filepath.Join("test", templateSpecPath)
		oldTemplateSpec = `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[external],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,true,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
`
		allInput.SharedInput.DestinationDirPrefix = ""
		allInput.ServiceInput.DirectoryName = destinationPath
		allInput.ServiceInput.TemplateName = "templateFolderStartPath"
		allInput.GeneratorInputs = map[string]string{
			"templateFolderStartPath": templateFolderStartPath,
		}
		allInput.ExtraInputs = map[string]string{
			"testInput":  "test worked",
			"testInput2": "test succeeded",
		}
		allInput.User = "external"
		testOptions.generationType = serviceStr
	})
	AfterEach(func() {
		os.RemoveAll(templateFolderStartPath)
		os.RemoveAll(destinationPath)
	})

	It("Generates template spec file when it doesn't exist", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(oldTemplateSpec))
	})

	It("Template Spec remains the same when spec doesn't change", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(oldTemplateSpec))
	})

	It("Required files are removed from spec when folders no longer have them", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		os.RemoveAll("test/test2")
		os.Remove("test/path0.go")
		os.Remove("test/test3/path5.go")
		newTemplateSpec := `File name,Overwrite,Users,Rename
test1,true,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[external],
test3,true,[external],
test3/path4.go,true,[external],
test3/path6.go,true,[external],
`
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(newTemplateSpec))
	})

	It("Maintain settings from old spec when files are added and removed", func() {
		templateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,false,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[msft],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,false,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
`
		// Generate template spec to be considered old spec
		err := stateFiles.WriteToFile(completeTemplateSpecPath, templateSpec)
		Expect(err).NotTo(HaveOccurred())
		os.Remove("test/test3/path5.go")
		os.Mkdir("test/test4", 0777)
		_, _ = os.Create("test/test4/path7.go")
		_, _ = os.Create("test/test4/path8.go")
		newTemplateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,false,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[msft],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,false,[external],
test3/path6.go,true,[external],
test4,true,[external],
test4/path7.go,true,[external],
test4/path8.go,true,[external],
`
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(newTemplateSpec))
	})

	It("Do not overwrite files that don't have overwrite capabilities", func() {
		oldTemplateSpec = `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,false,[external],
test1/path2.go,true,[msft],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,false,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
`
		// Generate template spec to be considered old spec
		err := stateFiles.WriteToFile(completeTemplateSpecPath, oldTemplateSpec)
		Expect(err).NotTo(HaveOccurred())
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		// Write to two of the files
		b := []byte("TEST WORKED!")
		err = ioutil.WriteFile("dest/test3/path4.go", b, 0644)
		Expect(err).NotTo(HaveOccurred())
		if err != nil {
			panic(err)
		}
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		// Upon regeneration file contents for the non overwrite should not have changed
		contents, err := ioutil.ReadFile("dest/test3/path4.go")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal(string(b)))
	})

	It("Do not overwrite files added by the user", func() {
		oldTemplateSpec = `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,false,[external],
test1/path2.go,true,[msft],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,false,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
`
		// Generate template spec to be considered old spec
		err := stateFiles.WriteToFile(completeTemplateSpecPath, oldTemplateSpec)
		Expect(err).NotTo(HaveOccurred())
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		// User creates a new  files and writes to one of the files
		_, _ = os.Create("dest/test3/path7.go")
		_, _ = os.Create("dest/test3/path8.go")
		b := []byte("TEST WORKED!")
		err = ioutil.WriteFile("dest/test3/path7.go", b, 0644)
		Expect(err).NotTo(HaveOccurred())

		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		// Upon regeneration files should exist and contents should not have changed
		Expect("dest/test3/path7.go").To(BeAnExistingFile())
		Expect("dest/test3/path8.go").To(BeAnExistingFile())
		contents, err := ioutil.ReadFile("dest/test3/path7.go")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal(string(b)))
	})

	It("Only generate files specific to provided user", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(oldTemplateSpec))
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		Expect("dest/path0.go").To(BeAnExistingFile())
		Expect("dest/test1").To(BeADirectory())
		Expect("dest/test1/path1.go").To(BeAnExistingFile())
		Expect("dest/test1/path2.go").To(BeAnExistingFile())
		Expect("dest/test2").To(BeADirectory())
		Expect("dest/test2/path3.go").To(BeAnExistingFile())
		Expect("dest/test3").To(BeADirectory())
		Expect("dest/test3/path4.go").To(BeAnExistingFile())
		Expect("dest/test3/path5.go").To(BeAnExistingFile())
		Expect("dest/test3/path6.go").To(BeAnExistingFile())
		templateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external msft],
test1/path1.go,true,[external],
test1/path2.go,true,[external msft],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external msft],
test3/path4.go,true,[external msft],
test3/path5.go,true,[external],
test3/path6.go,true,[external msft],
`
		newState := `test1
test1/path2.go
test3
test3/path4.go
test3/path6.go
`
		// Edit template spec to be new spec
		err = stateFiles.WriteToFile(completeTemplateSpecPath, templateSpec)
		Expect(err).NotTo(HaveOccurred())
		testOptions.user = "msft"
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		// Make sure only files for msft user exist
		Expect("dest/path0.go").ToNot(BeAnExistingFile())
		Expect("dest/test1").To(BeADirectory())
		Expect("dest/test1/path1.go").ToNot(BeAnExistingFile())
		Expect("dest/test1/path2.go").To(BeAnExistingFile())
		Expect("dest/test2").ToNot(BeADirectory())
		Expect("dest/test2/path3.go").ToNot(BeAnExistingFile())
		Expect("dest/test3").To(BeADirectory())
		Expect("dest/test3/path4.go").To(BeAnExistingFile())
		Expect("dest/test3/path5.go").ToNot(BeAnExistingFile())
		Expect("dest/test3/path6.go").To(BeAnExistingFile())
		//Check new state matches
		state, err := ioutil.ReadFile(filepath.Join(destinationPath, statePath))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(state)).To(Equal(newState))
	})

	It("Nested templating with one instance works", func() {
		newFolderPath := "test/test4"
		newTemplatesPath := filepath.Join(newFolderPath, nestedTemplatesPath)
		os.Mkdir(newFolderPath, 0777)
		os.Mkdir(newTemplatesPath, 0777)
		nestedJson :=
			`{
	"description": "",
	"id": "<<.extraInputs.testInput>>",
	"layout": {
		"height": 10,
		"width": 10,
		"x": 0,
		"y": 0
	}
}`
		parentJson :=
			`{
	"nested": <<template "nested.json" . >>,
	"testValue": "<<.extraInputs.testInput>>"
}`
		jsonPath := filepath.Join(newTemplatesPath, "nested.json")
		err := stateFiles.WriteToFile(jsonPath, nestedJson)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile("test/test4/parent.json", parentJson)
		Expect(err).NotTo(HaveOccurred())
		templateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[external],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,true,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
test4,true,[external],
test4/parent.json,true,[external],
`
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(templateSpec))
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		checkPath := filepath.Join("dest/test4", nestedTemplatesPath)
		Expect(checkPath).ToNot(BeADirectory())
		nestedOutput :=
			`{
	"nested": {
	"description": "",
	"id": "test worked",
	"layout": {
		"height": 10,
		"width": 10,
		"x": 0,
		"y": 0
	}
},
	"testValue": "test worked"
}`
		nestedVal, err := ioutil.ReadFile("dest/test4/parent.json")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(nestedVal)).To(Equal(nestedOutput))
	})

	It("Nested templating with multiple instances works", func() {
		newFolderPath1 := "test/test4"
		newTemplatesPath1 := filepath.Join(newFolderPath1, nestedTemplatesPath)
		os.Mkdir(newFolderPath1, 0777)
		os.Mkdir(newTemplatesPath1, 0777)
		newFolderPath2 := "test/test5"
		newTemplatesPath2 := filepath.Join(newFolderPath2, nestedTemplatesPath)
		os.Mkdir(newFolderPath2, 0777)
		os.Mkdir(newTemplatesPath2, 0777)
		nestedJson :=
			`{
	"description": "",
	"id": "<<.extraInputs.testInput>>",
	"layout": {
		"height": 10,
		"width": 10,
		"x": 0,
		"y": 0
	}
}`
		parentJson :=
			`{
	"nested": <<template "nested.json" . >>,
	"testValue": "<<.extraInputs.testInput>>"
}`
		nested1Json :=
			`{
	"description": "test",
	"id": "<<.extraInputs.testInput>>",
}`
		nested2Json :=
			`{
	"layout": {
		"height": "<<.extraInputs.testInput>>",
		"width": "<<.extraInputs.testInput2>>",
		"x": 0,
		"y": 0
	}
}`
		parent1Json :=
			`{
	"nested1": <<template "nested1.json" . >>,
	"nested2": <<template "nested2.json" . >>,
	"testValue": "<<.extraInputs.testInput>>"
}`
		parent2Json :=
			`{
	"nested2": <<template "nested2.json" . >>,
	"testValue": "<<.extraInputs.testInput>>"
}`
		parent3Json :=
			`{
	"test": "No variables"
}`
		nestedJsonPath := filepath.Join(newTemplatesPath1, "nested.json")
		err := stateFiles.WriteToFile(nestedJsonPath, nestedJson)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile("test/test4/parent.json", parentJson)
		Expect(err).NotTo(HaveOccurred())
		nestedJsonPath1 := filepath.Join(newTemplatesPath2, "nested1.json")
		nestedJsonPath2 := filepath.Join(newTemplatesPath2, "nested2.json")
		err = stateFiles.WriteToFile(nestedJsonPath1, nested1Json)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile(nestedJsonPath2, nested2Json)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile("test/test5/parent1.json", parent1Json)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile("test/test5/parent2.json", parent2Json)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile("test/test5/parent3.json", parent3Json)
		Expect(err).NotTo(HaveOccurred())
		templateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[external],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,true,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
test4,true,[external],
test4/parent.json,true,[external],
test5,true,[external],
test5/parent1.json,true,[external],
test5/parent2.json,true,[external],
test5/parent3.json,true,[external],
`
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(templateSpec))
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		checkPath1 := filepath.Join("dest/test4", nestedTemplatesPath)
		checkPath2 := filepath.Join("dest/test5", nestedTemplatesPath)
		Expect(checkPath1).ToNot(BeADirectory())
		Expect(checkPath2).ToNot(BeADirectory())
		nestedOutput :=
			`{
	"nested": {
	"description": "",
	"id": "test worked",
	"layout": {
		"height": 10,
		"width": 10,
		"x": 0,
		"y": 0
	}
},
	"testValue": "test worked"
}`
		nestedVal, err := ioutil.ReadFile("dest/test4/parent.json")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(nestedVal)).To(Equal(nestedOutput))
		nestedOutput1 :=
			`{
	"nested1": {
	"description": "test",
	"id": "test worked",
},
	"nested2": {
	"layout": {
		"height": "test worked",
		"width": "test succeeded",
		"x": 0,
		"y": 0
	}
},
	"testValue": "test worked"
}`
		nestedOutput2 :=
			`{
	"nested2": {
	"layout": {
		"height": "test worked",
		"width": "test succeeded",
		"x": 0,
		"y": 0
	}
},
	"testValue": "test worked"
}`
		nestedOutput3 :=
			`{
	"test": "No variables"
}`
		nestedVal, err = ioutil.ReadFile("dest/test5/parent1.json")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(nestedVal)).To(Equal(nestedOutput1))
		nestedVal, err = ioutil.ReadFile("dest/test5/parent2.json")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(nestedVal)).To(Equal(nestedOutput2))
		nestedVal, err = ioutil.ReadFile("dest/test5/parent3.json")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(nestedVal)).To(Equal(nestedOutput3))
	})

	It("Shell file is created with necessary permissions", func() {
		err := stateFiles.WriteToFile(completeTemplateSpecPath, oldTemplateSpec)
		Expect(err).NotTo(HaveOccurred())
		_, _ = os.Create("test/test3/shell.sh")
		shellScript := `echo "It Worked!"`
		err = stateFiles.WriteToFile("test/test3/shell.sh", shellScript)
		Expect(err).NotTo(HaveOccurred())
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())

		err = genFromTemplate(allInput, testOptions)
		Expect(err).ToNot(HaveOccurred())
		shell, err := ioutil.ReadFile("dest/test3/shell.sh")
		Expect(string(shell)).To(Equal(shellScript))
		Expect(err).NotTo(HaveOccurred())
		cmd := exec.Command("bash", "dest/test3/shell.sh")
		stdout, err := cmd.Output()
		Expect(string(stdout)).To(Equal("It Worked!\n"))

	})

	It("Temporary folder doesn't exist after succesful generation", func() {
		oldTemplateSpec = `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,false,[external],
test1/path2.go,true,[msft],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,false,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
`
		// Generate template spec to be considered old spec
		err := stateFiles.WriteToFile(completeTemplateSpecPath, oldTemplateSpec)
		Expect(err).NotTo(HaveOccurred())
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())

		Expect("dest/.temp").ToNot(BeADirectory())
	})

	It("Errors when state file doesn't exist after creating once", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())

		err = genFromTemplate(allInput, testOptions)
		Expect(err).ToNot(HaveOccurred())

		os.Remove(filepath.Join(destinationPath, statePath))
		err = genFromTemplate(allInput, testOptions)
		Expect(err).To(HaveOccurred())

	})

	It("Errors when template spec includes items that template folder doesn't have", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		// Create fake testOptions for generation

		os.Remove("test/test2/path3.go")
		err = genFromTemplate(allInput, testOptions)
		Expect(err).To(HaveOccurred())
		os.RemoveAll(destinationPath)
		_, _ = os.Create("test/test2/path3.go")
		err = genFromTemplate(allInput, testOptions)
		Expect(err).NotTo(HaveOccurred())
		os.RemoveAll(destinationPath)
		os.RemoveAll("test/test2")
		err = genFromTemplate(allInput, testOptions)
		Expect(err).To(HaveOccurred())
	})

	It("Error handling incorrect template contents for parsing", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(oldTemplateSpec))
		testFile :=
			`"testValue": "<<.extraInputs.testInput"`
		err = stateFiles.WriteToFile("test/test2/path3.go", testFile)
		Expect(err).NotTo(HaveOccurred())
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).To(HaveOccurred())
		Expect("dest/.temp").ToNot(BeADirectory())
	})

	It("Error handling with missing nested template parsefiles", func() {
		newFolderPath := "test/test4"
		newTemplatesPath := filepath.Join(newFolderPath, nestedTemplatesPath)
		os.Mkdir(newFolderPath, 0777)
		os.Mkdir(newTemplatesPath, 0777)
		parentJson :=
			`{
	"nested": <<template "missingNested.json" . >>,
	"testValue": "<<.extraInputs.testInput>>"
}`
		err := stateFiles.WriteToFile("test/test4/parent.json", parentJson)
		Expect(err).NotTo(HaveOccurred())
		templateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[external],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,true,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
test4,true,[external],
test4/parent.json,true,[external],
`
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(templateSpec))
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)

		Expect(err).To(HaveOccurred())
		Expect("dest/.temp").ToNot(BeADirectory())
	})

	It("Error handling with incorrect template execution", func() {
		newFolderPath := "test/test4"
		newTemplatesPath := filepath.Join(newFolderPath, nestedTemplatesPath)
		os.Mkdir(newFolderPath, 0777)
		os.Mkdir(newTemplatesPath, 0777)
		nestedJson :=
			`{
	"description": "",
	"id": "<<.wrongVar>>",
	"layout": {
		"height": 10,
		"width": 10,
		"x": 0,
		"y": 0
	}
}`
		parentJson :=
			`{
	"nested": <<template "nested.json" . >>,
	"testValue": "<<.testInput>>"
}`
		jsonPath := filepath.Join(newTemplatesPath, "nested.json")
		err := stateFiles.WriteToFile(jsonPath, nestedJson)
		Expect(err).NotTo(HaveOccurred())
		err = stateFiles.WriteToFile("test/test4/parent.json", parentJson)
		Expect(err).NotTo(HaveOccurred())
		templateSpec := `File name,Overwrite,Users,Rename
path0.go,true,[external],
test1,true,[external],
test1/path1.go,true,[external],
test1/path2.go,true,[external],
test2,true,[external],
test2/path3.go,true,[external],
test3,true,[external],
test3/path4.go,true,[external],
test3/path5.go,true,[external],
test3/path6.go,true,[external],
test4,true,[external],
test4/parent.json,true,[external],
`
		err = generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(templateSpec))
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)

		Expect(err).To(HaveOccurred())
		Expect("dest/.temp").ToNot(BeADirectory())
	})

	It("Error handling destination path being invalid", func() {

		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(oldTemplateSpec))

		allInput.ServiceInput.DirectoryName = strings.Repeat("a", 4096)
		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).To(HaveOccurred())
		Expect("dest/.temp").ToNot(BeADirectory())
		os.RemoveAll(strings.Repeat("a", 4096))
	})

	It("Error handling copying file to destination", func() {
		err := generateTemplateSpec(templateFolderStartPath, testOptions)
		Expect(err).NotTo(HaveOccurred())
		spec, err := ioutil.ReadFile(completeTemplateSpecPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(spec)).To(Equal(oldTemplateSpec))

		// Create fake testOptions for generation
		err = genFromTemplate(allInput, testOptions)
		Expect(err).ToNot(HaveOccurred())
		Expect("dest/.temp").ToNot(BeADirectory())

		os.Remove("dest/test3/path6.go")

		// We remove all permissions from the file
		origUmask := syscall.Umask(0)
		_, err = os.OpenFile("dest/test3/path6.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o000)
		Expect(err).ToNot(HaveOccurred())
		syscall.Umask(origUmask)

		// Upon regeneration error will occur in copying file as we have removed permissions
		err = genFromTemplate(allInput, testOptions)
		Expect(err).To(HaveOccurred())
		Expect("dest/.temp").ToNot(BeADirectory())
	})
})
