package images_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestImages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Images Suite")
}
