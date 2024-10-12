package notifier_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSecretsanta(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Notifier Suite")
}
