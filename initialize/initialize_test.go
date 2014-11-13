package initialize_test

import (
	. "github.com/it0a/last98/initialize"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type stubEnvReader struct {
	port string
}

func (s stubEnvReader) ReadPort() string {
	return s.port
}

var _ = Describe("Init", func() {

	Describe("reading the $PORT env var", func() {

		Context("when it is not set", func() {
			var envReader stubEnvReader
			envReader.port = ""
			It("should default to 8080", func() {
				Expect(ReadPort(envReader)).To(Equal("8080"))
			})
		})

		Context("when it is set", func() {
			var envReader stubEnvReader
			envReader.port = "1234"
			It("should read in the expected value", func() {
				Expect(ReadPort(envReader)).To(Equal("1234"))
			})
		})

	})
})
