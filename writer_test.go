package cdb64_test

import (
	"github.com/bsm/cdb64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", func() {
	var subject *cdb64.Writer

	BeforeEach(func() {
		w, err := cdb64.Create(testDir + "/test.cdb")
		Expect(err).NotTo(HaveOccurred())
		subject = w
	})

	AfterEach(func() {
		Expect(subject.Close()).To(Succeed())
	})

	It("should PUT", func() {
		Expect(seedData(subject, 1000)).To(Succeed())
	})
})
