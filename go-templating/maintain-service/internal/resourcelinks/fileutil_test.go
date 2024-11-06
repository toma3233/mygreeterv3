package resourcelinks

import (
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceLinks", func() {
	dest := "test"
	svcPath := filepath.Join(dest, "mysvc")
	svcFileName := "svc_resources.md"
	genSvcFilePath := filepath.Join(svcPath, svcFileName)
	resourceDir1 := filepath.Join(svcPath, "resourceDir1")
	resourceDir2 := filepath.Join(svcPath, "resourceDir2")
	nonResourceDir := filepath.Join(svcPath, "nonResourceDir")
	testdataDir := "testdata"
	testDeployFileContent, err := os.ReadFile(filepath.Join(testdataDir, "deploy_resources.md"))
	if err != nil {
		log.Printf("Error reading file: %v", err)
	}
	testSvcFileContent, err := os.ReadFile(filepath.Join(testdataDir, "svc_resources.md"))
	if err != nil {
		log.Printf("Error reading file: %v", err)
	}
	BeforeEach(func() {
		os.Mkdir(dest, 0777)
		os.Mkdir(svcPath, 0777)
		os.Mkdir(resourceDir1, 0777)
		os.Mkdir(resourceDir2, 0777)
		os.Mkdir(nonResourceDir, 0777)
		jsonFileContent, err := os.ReadFile(filepath.Join(testdataDir, ".deploy_output.json"))
		if err != nil {
			log.Printf("Error reading file: %v", err)
		}
		err = os.WriteFile(filepath.Join(resourceDir1, ".resource1_output.json"), jsonFileContent, 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}
		err = os.WriteFile(filepath.Join(resourceDir1, ".resource2_output.json"), jsonFileContent, 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}
		err = os.WriteFile(filepath.Join(resourceDir2, ".resource1_output.json"), jsonFileContent, 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}

	})
	AfterEach(func() {
		os.RemoveAll(dest)
	})
	It("Generates service level markdown file and all deployment level markdown files when desired json files are valid.", func() {
		err := GenerateAllResourceFiles(svcPath, svcFileName)
		Expect(err).ToNot(HaveOccurred())

		filePath1 := filepath.Join(resourceDir1, "resource1_resources.md")
		Expect(filePath1).To(BeAnExistingFile())
		filePath2 := filepath.Join(resourceDir1, "resource2_resources.md")
		Expect(filePath2).To(BeAnExistingFile())
		filePath3 := filepath.Join(resourceDir2, "resource1_resources.md")
		Expect(filePath3).To(BeAnExistingFile())
		nonFilePath4 := filepath.Join(nonResourceDir, "nonResource_resources.md")
		Expect(nonFilePath4).ToNot(BeAnExistingFile())
		Expect(genSvcFilePath).To(BeAnExistingFile())

		genDeployFileContent, err := os.ReadFile(filePath1)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(genDeployFileContent)).To(Equal(string(testDeployFileContent)))

		genSvcFileContent, err := os.ReadFile(genSvcFilePath)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(genSvcFileContent)).To(Equal(string(testSvcFileContent)))

	})
	It("Deploy level markdown file not created when the respective json file is invalid.", func() {
		invalidFileContent := `Invalid json file`
		err := os.WriteFile(filepath.Join(resourceDir1, ".invalid_output.json"), []byte(invalidFileContent), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = GenerateAllResourceFiles(svcPath, svcFileName)
		Expect(err).ToNot(HaveOccurred())

		Expect(filepath.Join(resourceDir1, "invalid_markdown.md")).ToNot(BeAnExistingFile())

		genSvcFileContent, err := os.ReadFile(genSvcFilePath)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(genSvcFileContent)).To(Equal(string(testSvcFileContent)))
	})
	It("No service level markdown file created if all deployment level files in same directory.", func() {
		err := GenerateAllResourceFiles(resourceDir1, svcFileName)
		Expect(err).ToNot(HaveOccurred())
		filePath1 := filepath.Join(resourceDir1, svcFileName)
		Expect(filePath1).ToNot(BeAnExistingFile())
	})
	It("Service path directory doesn't exist.", func() {
		err := GenerateAllResourceFiles("incorrectSvcPath", svcFileName)
		Expect(err).To(HaveOccurred())
	})
})
