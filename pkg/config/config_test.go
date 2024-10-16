package config_test

import (
	"github.com/igodwin/secretsanta/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"text/template"
)

const (
	expectedSMTPHost        = "smtp.example.com"
	expectedSMTPPort        = "587"
	expectedSMTPIdentity    = ""
	expectedSMTPUsername    = "user@example.com"
	expectedSMTPPassword    = "password"
	expectedSMTPFromAddress = "no-reply@example.com"
	expectedSMTPFromName    = ""

	configTemplate = `[smtp]
host = "{{.Host}}"
port = "{{.Port}}"
identity = "{{.Identity}}"
username = "{{.Username}}"
password = "{{.Password}}"
from_address = "{{.FromAddress}}"
from_name = "{{.FromName}}"
`
)

var _ = Describe("Config", func() {
	var tempDir string
	var testConfig *config.Config

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "config_test")
		Expect(err).NotTo(HaveOccurred())

		smtp := config.SMTPConfig{
			Host:        expectedSMTPHost,
			Port:        expectedSMTPPort,
			Identity:    expectedSMTPIdentity,
			Username:    expectedSMTPUsername,
			Password:    expectedSMTPPassword,
			FromAddress: expectedSMTPFromAddress,
			FromName:    expectedSMTPFromName,
		}
		configFilePath := filepath.Join(tempDir, "secretsanta.config")
		file, err := os.Create(configFilePath)
		Expect(err).NotTo(HaveOccurred())

		t := template.Must(template.New("validConfig").Parse(configTemplate))
		Expect(t.Execute(file, smtp)).NotTo(HaveOccurred())

		Expect(file.Close()).To(Succeed())

		config.Paths = []string{tempDir}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	Context("using valid config", func() {
		Describe("GetConfig", func() {
			It("should return initialized Config with expected values", func() {
				testConfig = config.GetConfig()
				Expect(testConfig).NotTo(BeNil())
				Expect(testConfig.SMTP).NotTo(BeNil())
				Expect(testConfig.SMTP.Host).To(Equal(expectedSMTPHost))
				Expect(testConfig.SMTP.Port).To(Equal(expectedSMTPPort))
				Expect(testConfig.SMTP.Identity).To(Equal(expectedSMTPIdentity))
				Expect(testConfig.SMTP.Username).To(Equal(expectedSMTPUsername))
				Expect(testConfig.SMTP.Password).To(Equal(expectedSMTPPassword))
				Expect(testConfig.SMTP.FromAddress).To(Equal(expectedSMTPFromAddress))
				Expect(testConfig.SMTP.FromName).To(Equal(expectedSMTPFromName))
			})

			It("should return same pointer for subsequent calls", func() {
				testConfig = config.GetConfig()
				testConfig2 := config.GetConfig()
				testConfig3 := config.GetConfig()

				Expect(testConfig).To(BeIdenticalTo(testConfig2))
				Expect(testConfig).To(BeIdenticalTo(testConfig3))
			})
		})
	})
})
