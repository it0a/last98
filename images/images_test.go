package images_test

import (
	. "github.com/it0a/last98/images"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Images", func() {
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
