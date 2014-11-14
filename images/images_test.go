package images_test

import (
	"errors"
	. "github.com/it0a/last98/images"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type StubImageDatabase struct {
	id        string
	imageData ImageData
	err       error
}

// We'll only return images if the id is equal to 1
func (stub StubImageDatabase) FindById(id string) (ImageData, error) {
	if id == "1" {
		return stub.imageData, nil
	} else {
		return ImageData{}, errors.New("error")
	}
}

func (stub StubImageDatabase) Save(i NewImageData) error {
	return nil
}

func (stub StubImageDatabase) Delete(id string) error {
	if id != "1" {
		return errors.New("Error!")
	}
	return nil
}

var _ = Describe("Images", func() {

	Describe("image database", func() {

		expected := "expected"
		stubImageDatabase := StubImageDatabase{id: "1", imageData: ImageData{Data: expected}}

		Describe("reading from a repository", func() {

			Context("we retrieve an ID that exists", func() {
				It("should receive valid image data", func() {
					imageData, _ := ReadImage("1", stubImageDatabase)
					Expect(imageData.Data).To(Equal(expected))
				})
			})

			Context("when reading a non-existant ID", func() {

				It("should not receive valid image data", func() {
					imageData, _ := ReadImage("100", stubImageDatabase)
					Expect(imageData).ShouldNot(Equal(expected))
				})

				It("should have received an error", func() {
					_, err := ReadImage("100", stubImageDatabase)
					Expect(err).Should(HaveOccurred())
				})

			})
		})

		Describe("image creation", func() {
			Context("when saving an image sucessfully", func() {
				It("returns with no errors", func() {
					err := SaveImage(NewImageData{}, stubImageDatabase)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		Describe("image deletion", func() {
			Context("We delete an image", func() {
				It("returns with no errors", func() {
					err := DeleteImage("1", stubImageDatabase)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("We unsuccessfully delete an image", func() {
				It("causes an error", func() {
					err := DeleteImage("2", stubImageDatabase)
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
