package populatemethods

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PopulateMethods", func() {
	dest := "test"
	protoPath := filepath.Join(dest, "api.proto")
	svcPath := filepath.Join(dest, "mysvc")
	deleteGarbageFiles := false
	templDir1 := filepath.Join(svcPath, "templateDir1")
	templDir2 := filepath.Join(svcPath, "templateDir2")
	templDir3 := filepath.Join(svcPath, "templateDir3")
	BeforeEach(func() {
		os.Mkdir(dest, 0777)
		os.Mkdir(svcPath, 0777)
		os.Mkdir(templDir1, 0777)
		os.Mkdir(templDir2, 0777)
		os.Mkdir(templDir3, 0777)
		_, _ = os.Create("test/api.proto")
		// directory contains method template, method and method state file
		_, _ = os.Create("test/mysvc/templateDir1/.method_template_go.txt")
		_, _ = os.Create("test/mysvc/templateDir1/SayHello.go")
		_, _ = os.Create("test/mysvc/templateDir1/.methods_state.txt")
		// directory contains method template, method and method state file
		_, _ = os.Create("test/mysvc/templateDir2/.method_template_go.txt")
		_, _ = os.Create("test/mysvc/templateDir2/SayHello.go")
		_, _ = os.Create("test/mysvc/templateDir2/.methods_state.txt")
		// directory contains only method template file
		_, _ = os.Create("test/mysvc/templateDir3/.method_template_json.txt")
		protoContent := `
		service mysvc {
		rpc SayHello (HelloRequest) returns (HelloReply) {}
		}
		`
		err := ioutil.WriteFile("test/api.proto", []byte(protoContent), 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}
		stateContent := "SayHello.go\n"
		err = ioutil.WriteFile("test/mysvc/templateDir1/.methods_state.txt", []byte(stateContent), 0644)
		err = ioutil.WriteFile("test/mysvc/templateDir2/.methods_state.txt", []byte(stateContent), 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}

	})
	AfterEach(func() {
		os.RemoveAll(dest)
	})
	It("No change to method files or states files when proto not modified.", func() {
		methodFilePath := filepath.Join(templDir1, "SayHello.go")
		methodStateFilePath := filepath.Join(templDir1, ".methods_state.txt")
		methodFileContent := "no change"

		err := ioutil.WriteFile(methodFilePath, []byte(methodFileContent), 0644)
		Expect(err).ToNot(HaveOccurred())

		err = PopulateAllMethodFiles(svcPath, protoPath, deleteGarbageFiles)
		Expect(err).ToNot(HaveOccurred())

		newMethodFileContent, err := ioutil.ReadFile(methodFilePath)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(newMethodFileContent)).To(Equal(methodFileContent))

		newStateContent, err := ioutil.ReadFile(methodStateFilePath)
		Expect(err).ToNot(HaveOccurred())
		stateContent := "SayHello.go\n"
		Expect(string(newStateContent)).To(Equal(stateContent))
	})
	It("Generates state file and method file when they don't exist.", func() {
		templPath := filepath.Join(templDir3, ".method_template_json.txt")
		newMethodFilePath := filepath.Join(templDir3, "SayHello.json")
		methodStateFilePath := filepath.Join(templDir3, ".methods_state.txt")

		Expect(templPath).To(BeAnExistingFile())
		Expect(newMethodFilePath).ToNot(BeAnExistingFile())
		Expect(methodStateFilePath).ToNot(BeAnExistingFile())

		err := PopulateAllMethodFiles(svcPath, protoPath, deleteGarbageFiles)
		Expect(err).NotTo(HaveOccurred())

		Expect(newMethodFilePath).To(BeAnExistingFile())
		Expect(methodStateFilePath).To(BeAnExistingFile())

		stateContent := "SayHello.json\n"
		state, err := ioutil.ReadFile(methodStateFilePath)
		Expect(string(state)).To(Equal(stateContent))
	})
	It("Proto method addition triggers method file creation and state file updates in template directories.", func() {
		newProtoContent := `
		service mysvc {
		rpc SayHello (HelloRequest) returns (HelloReply) {}
		rpc SayGoodbye (GoodbyeRequest) returns (GoodbyeReply) {}
		}
		`
		err := ioutil.WriteFile(protoPath, []byte(newProtoContent), 0644)
		Expect(err).NotTo(HaveOccurred())

		invalidTemplContent := "{{.Name}}"
		err = ioutil.WriteFile("test/mysvc/templateDir1/.method_template_go.txt", []byte(invalidTemplContent), 0644)
		Expect(err).NotTo(HaveOccurred())

		existingMethodFilePath1 := filepath.Join(templDir1, "SayHello.go")
		newMethodFilePath1 := filepath.Join(templDir1, "SayGoodbye.go")
		methodStateFilePath1 := filepath.Join(templDir1, ".methods_state.txt")

		existingMethodFilePath2 := filepath.Join(templDir2, "SayHello.go")
		newMethodFilePath2 := filepath.Join(templDir2, "SayGoodbye.go")
		methodStateFilePath2 := filepath.Join(templDir2, ".methods_state.txt")

		err = PopulateAllMethodFiles(svcPath, protoPath, deleteGarbageFiles)
		Expect(err).ToNot(HaveOccurred())

		Expect(newMethodFilePath1).To(BeAnExistingFile())
		Expect(newMethodFilePath2).To(BeAnExistingFile())
		Expect(existingMethodFilePath1).To(BeAnExistingFile())
		Expect(existingMethodFilePath2).To(BeAnExistingFile())

		stateContent := "SayHello.go\nSayGoodbye.go\n"
		state, err := ioutil.ReadFile(methodStateFilePath1)
		Expect(string(state)).To(Equal(stateContent))
		state, err = ioutil.ReadFile(methodStateFilePath2)
		Expect(string(state)).To(Equal(stateContent))

		fileContent := "SayGoodbye"
		generatedFileContent, err := ioutil.ReadFile(newMethodFilePath1)
		Expect(string(generatedFileContent)).To(Equal(fileContent))
	})
	It("Deletion of method in proto file updates method state file and doesn't delete respective method files.", func() {
		newProtoContent := `
		service mysvc {
		}
		`
		err := ioutil.WriteFile(protoPath, []byte(newProtoContent), 0644)
		Expect(err).ToNot(HaveOccurred())

		existingMethodFilePath1 := filepath.Join(templDir1, "SayHello.go")
		methodStateFilePath1 := filepath.Join(templDir1, ".methods_state.txt")
		existingMethodFilePath2 := filepath.Join(templDir2, "SayHello.go")
		methodStateFilePath2 := filepath.Join(templDir2, ".methods_state.txt")

		err = PopulateAllMethodFiles(svcPath, protoPath, deleteGarbageFiles)
		Expect(err).ToNot(HaveOccurred())

		Expect(existingMethodFilePath1).To(BeAnExistingFile())
		Expect(existingMethodFilePath2).To(BeAnExistingFile())

		stateContent := ""
		state, err := ioutil.ReadFile(methodStateFilePath1)
		Expect(string(state)).To(Equal(stateContent))
		state, err = ioutil.ReadFile(methodStateFilePath2)
		Expect(string(state)).To(Equal(stateContent))
	})
	It("Deletion of method in proto file updates method state file and deletes respective method files.", func() {
		deleteGarbageFiles = true // Automatically delete files

		newProtoContent := `
		service mysvc {
		}
		`

		err := ioutil.WriteFile(protoPath, []byte(newProtoContent), 0644)
		Expect(err).ToNot(HaveOccurred())

		existingMethodFilePath1 := filepath.Join(templDir1, "SayHello.go")
		methodStateFilePath1 := filepath.Join(templDir1, ".methods_state.txt")
		existingMethodFilePath2 := filepath.Join(templDir2, "SayHello.go")
		methodStateFilePath2 := filepath.Join(templDir2, ".methods_state.txt")

		err = PopulateAllMethodFiles(svcPath, protoPath, deleteGarbageFiles)
		Expect(err).ToNot(HaveOccurred())

		Expect(existingMethodFilePath1).ToNot(BeAnExistingFile())
		Expect(existingMethodFilePath2).ToNot(BeAnExistingFile())

		stateContent := ""
		state, err := ioutil.ReadFile(methodStateFilePath1)
		Expect(string(state)).To(Equal(stateContent))
		state, err = ioutil.ReadFile(methodStateFilePath2)
		Expect(string(state)).To(Equal(stateContent))
	})
	It("With an invalid template, the file isn’t created but traversal persists.", func() {
		newProtoContent := `
		service mysvc {
		rpc SayHello (HelloRequest) returns (HelloReply) {}
		rpc SayGoodbye (GoodbyeRequest) returns (GoodbyeReply) {}
		}
		`
		err := ioutil.WriteFile(protoPath, []byte(newProtoContent), 0644)
		invalidTemplContent := "{{.Name}\n"
		err = ioutil.WriteFile("test/mysvc/templateDir1/.method_template_go.txt", []byte(invalidTemplContent), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = PopulateAllMethodFiles(svcPath, protoPath, deleteGarbageFiles)
		Expect(err).ToNot(HaveOccurred())
		fileName := filepath.Join(templDir1, "SayGoodbye.go")
		Expect(fileName).ToNot(BeAnExistingFile())
	})
	It("Api proto file path doesn't exist", func() {
		err := PopulateAllMethodFiles(svcPath, "incorrectProtoPath.proto", deleteGarbageFiles)
		Expect(err).To(HaveOccurred())
	})
	It("Service path directory doesn't exist", func() {
		err := PopulateAllMethodFiles("incorrectSvcPath", protoPath, deleteGarbageFiles)
		Expect(err).To(HaveOccurred())
	})
})
