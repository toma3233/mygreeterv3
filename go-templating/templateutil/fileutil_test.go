package templateutil

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shared functions", func() {
	folderPath := "test"
	outputFilePath := filepath.Join(folderPath, "output.txt")
	tmplPath := filepath.Join(folderPath, "template.txt")
	templatePrefix := "template-"
	envConfig := filepath.Join(folderPath, "env-config.yaml")
	envConfigContents := `resourcesName: TGEb373NDM
subscriptionId: 18134071`
	delimiter := []string{"{{", "}}"}
	BeforeEach(func() {
		os.Mkdir(folderPath, 0777)
		_, _ = os.Create(tmplPath)
		err := os.WriteFile(envConfig, []byte(envConfigContents), 0644)
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		os.RemoveAll(folderPath)
	})
	It("Output file created and populated when template is valid.", func() {
		tmplContent := "{{.Name}}"
		data := struct {
			Name string
		}{
			Name: "test",
		}
		err := os.WriteFile(tmplPath, []byte(tmplContent), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = CreateFileFromTemplate(outputFilePath, tmplPath, data, []string{"{{", "}}"}, 0777)
		Expect(err).ToNot(HaveOccurred())
		Expect(outputFilePath).To(BeAnExistingFile())
		fileContent := "test"
		newFileContent, err := os.ReadFile(outputFilePath)
		Expect(string(newFileContent)).To(Equal(fileContent))
	})
	It("Output file not created when error occurs", func() {
		tmplContent := "{{.Name}"
		data := struct {
			Name string
		}{
			Name: "test",
		}
		err := os.WriteFile(tmplPath, []byte(tmplContent), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = CreateFileFromTemplate(outputFilePath, tmplPath, data, []string{"{{", "}}"}, 0777)
		Expect(err).To(HaveOccurred())
		Expect(outputFilePath).ToNot(BeAnExistingFile())
	})
	It("All instances of template files are templated correctly.", func() {
		json1 :=
			`{
	"acrName": "servicehub-{{.resourcesName}}-acr",
	"serverSubscription": "{{.subscriptionId}}"
}`
		json2 :=
			`{
	"rbacName": "servicehub-{{.resourcesName}}-rbac",
	"serverSubscription": "{{.subscriptionId}}"
}`
		json3 :=
			`{
	"clusterName": "servicehub-{{.resourcesName}}-cluster",
	"serverSubscription": "{{.subscriptionId}}"
}`
		err := os.WriteFile(filepath.Join(folderPath, "template-file1.json"), []byte(json1), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(folderPath, "template-file2.json"), []byte(json2), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(folderPath, "template-file3.json"), []byte(json3), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = ExecuteTemplatesInFolder(templatePrefix, folderPath, envConfig, delimiter)
		Expect(err).ToNot(HaveOccurred())
		filePath1 := filepath.Join(folderPath, "file1.json")
		Expect(filePath1).To(BeAnExistingFile())
		filePath2 := filepath.Join(folderPath, "file2.json")
		Expect(filePath2).To(BeAnExistingFile())
		filePath3 := filepath.Join(folderPath, "file3.json")
		Expect(filePath3).To(BeAnExistingFile())
		templated1 :=
			`{
	"acrName": "servicehub-TGEb373NDM-acr",
	"serverSubscription": "18134071"
}`
		templated2 :=
			`{
	"rbacName": "servicehub-TGEb373NDM-rbac",
	"serverSubscription": "18134071"
}`
		templated3 :=
			`{
	"clusterName": "servicehub-TGEb373NDM-cluster",
	"serverSubscription": "18134071"
}`

		templated1FileContent, err := os.ReadFile(filePath1)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated1FileContent)).To(Equal(string(templated1)))

		templated2FileContent, err := os.ReadFile(filePath2)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated2FileContent)).To(Equal(string(templated2)))

		templated3FileContent, err := os.ReadFile(filePath3)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated3FileContent)).To(Equal(string(templated3)))
	})
	It("Different delimiters work as expected.", func() {
		json1 :=
			`{
	"acrName": "servicehub-<<.resourcesName>>-acr",
	"serverSubscription": "<<.subscriptionId>>"
}`
		json2 :=
			`{
	"rbacName": "servicehub-<<.resourcesName>>-rbac",
	"serverSubscription": "<<.subscriptionId>>"
}`
		json3 :=
			`{
	"clusterName": "servicehub-<<.resourcesName>>-cluster",
	"serverSubscription": "<<.subscriptionId>>"
}`
		err := os.WriteFile(filepath.Join(folderPath, "template-file1.json"), []byte(json1), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(folderPath, "template-file2.json"), []byte(json2), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(folderPath, "template-file3.json"), []byte(json3), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = ExecuteTemplatesInFolder(templatePrefix, folderPath, envConfig, []string{"<<", ">>"})
		Expect(err).ToNot(HaveOccurred())
		filePath1 := filepath.Join(folderPath, "file1.json")
		Expect(filePath1).To(BeAnExistingFile())
		filePath2 := filepath.Join(folderPath, "file2.json")
		Expect(filePath2).To(BeAnExistingFile())
		filePath3 := filepath.Join(folderPath, "file3.json")
		Expect(filePath3).To(BeAnExistingFile())
		templated1 :=
			`{
	"acrName": "servicehub-TGEb373NDM-acr",
	"serverSubscription": "18134071"
}`
		templated2 :=
			`{
	"rbacName": "servicehub-TGEb373NDM-rbac",
	"serverSubscription": "18134071"
}`
		templated3 :=
			`{
	"clusterName": "servicehub-TGEb373NDM-cluster",
	"serverSubscription": "18134071"
}`

		templated1FileContent, err := os.ReadFile(filePath1)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated1FileContent)).To(Equal(string(templated1)))

		templated2FileContent, err := os.ReadFile(filePath2)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated2FileContent)).To(Equal(string(templated2)))

		templated3FileContent, err := os.ReadFile(filePath3)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated3FileContent)).To(Equal(string(templated3)))
	})
	It("Extra variables not mention in the config file cause failure.", func() {
		json1 :=
			`{
	"acrName": "servicehub-{{.resourcesName}}-acr",
	"serverSubscription": "{{.subscriptionId}}"
	"random value": "{{.nonsenseVar}}"
}`
		err := os.WriteFile(filepath.Join(folderPath, "template-file1.json"), []byte(json1), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = ExecuteTemplatesInFolder(templatePrefix, folderPath, envConfig, delimiter)
		Expect(err).To(HaveOccurred())
		filePath1 := filepath.Join(folderPath, "file1.json")
		Expect(filePath1).ToNot(BeAnExistingFile())
	})
	It("Nested folders with template files function.", func() {
		json1 :=
			`{
	"acrName": "servicehub-{{.resourcesName}}-acr",
	"serverSubscription": "{{.subscriptionId}}"
}`
		json2 :=
			`{
	"rbacName": "servicehub-{{.resourcesName}}-rbac",
	"serverSubscription": "{{.subscriptionId}}"
}`
		json3 :=
			`{
	"clusterName": "servicehub-{{.resourcesName}}-cluster",
	"serverSubscription": "{{.subscriptionId}}"
}`
		innerFolderPath := filepath.Join(folderPath, "innerFolder")
		os.Mkdir(innerFolderPath, 0777)
		err := os.WriteFile(filepath.Join(folderPath, "template-file1.json"), []byte(json1), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(innerFolderPath, "template-file2.json"), []byte(json2), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(innerFolderPath, "template-file3.json"), []byte(json3), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = ExecuteTemplatesInFolder(templatePrefix, folderPath, envConfig, delimiter)
		Expect(err).ToNot(HaveOccurred())
		filePath1 := filepath.Join(folderPath, "file1.json")
		Expect(filePath1).To(BeAnExistingFile())
		filePath2 := filepath.Join(innerFolderPath, "file2.json")
		Expect(filePath2).To(BeAnExistingFile())
		filePath3 := filepath.Join(innerFolderPath, "file3.json")
		Expect(filePath3).To(BeAnExistingFile())
		templated1 :=
			`{
	"acrName": "servicehub-TGEb373NDM-acr",
	"serverSubscription": "18134071"
}`
		templated2 :=
			`{
	"rbacName": "servicehub-TGEb373NDM-rbac",
	"serverSubscription": "18134071"
}`
		templated3 :=
			`{
	"clusterName": "servicehub-TGEb373NDM-cluster",
	"serverSubscription": "18134071"
}`

		templated1FileContent, err := os.ReadFile(filePath1)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated1FileContent)).To(Equal(string(templated1)))

		templated2FileContent, err := os.ReadFile(filePath2)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated2FileContent)).To(Equal(string(templated2)))

		templated3FileContent, err := os.ReadFile(filePath3)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(templated3FileContent)).To(Equal(string(templated3)))
	})
	It("Folder path directory doesn't exist.", func() {
		err := ExecuteTemplatesInFolder(templatePrefix, "incorrectPath", envConfig, delimiter)
		Expect(err).To(HaveOccurred())
	})
})
