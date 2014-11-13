package images_test

import (
	"errors"
	. "github.com/it0a/last98/images"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type stubImageReader struct {
	id        string
	imageData ImageData
}

// We'll only return images if the id is equal to 1
func (stubImageReader stubImageReader) FindById(id string) (ImageData, error) {
	if id == "1" {
		return stubImageReader.imageData, nil
	} else {
		return ImageData{}, errors.New("error")
	}
}

var _ = Describe("Images", func() {

	Describe("image data", func() {

		Describe("reading from a repository", func() {

			// Setup
			expected := "expected"
			imageReader := stubImageReader{id: "1", imageData: ImageData{Data: expected}}

			Context("we retrieve an ID that exists", func() {
				It("should receive valid image data", func() {
					imageData, _ := ReadImage("1", imageReader)
					Expect(imageData.Data).To(Equal(expected))
				})
			})

			Context("when reading a non-existant ID", func() {

				It("should not receive valid image data", func() {
					imageData, _ := ReadImage("100", imageReader)
					Expect(imageData).ShouldNot(Equal(expected))
				})

				It("should have received an error", func() {
					_, err := ReadImage("100", imageReader)
					Expect(err).Should(HaveOccurred())
				})

			})
		})
	})

	Describe("Template Functions", func() {
		Describe("detecting when to end rendering a row of thumbnails", func() {

			Context("when rendering the first thumbnail", func() {
				It("should not be the end of a row", func() {
					Expect(IsEndOfRow(0)).To(Equal(false))
				})
			})

			Context("when rendering the fifth thumbnail", func() {
				It("should be the end of a row", func() {
					Expect(IsEndOfRow(4)).To(Equal(true))
				})
			})

		})
	})
})
