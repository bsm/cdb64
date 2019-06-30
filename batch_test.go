package cdb64_test

import (
	"github.com/bsm/cdb64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Batch", func() {
	var subject *cdb64.Batch
	var reader *cdb64.Reader

	BeforeEach(func() {
		w, err := cdb64.Create(testDir + "/test.cdb")
		Expect(err).NotTo(HaveOccurred())
		defer w.Close()

		Expect(seedData(w, 1000)).To(Succeed())

		reader, err = w.Freeze()
		Expect(err).NotTo(HaveOccurred())

		subject = reader.Batch()
	})

	AfterEach(func() {
		Expect(reader.Close()).To(Succeed())
	})

	It("should Get", func() {
		Expect(subject.Get([]byte("missing"))).To(BeNil())
		Expect(subject.Get([]byte("key-00000005"))).To(Equal([]byte("val-00000005")))
		Expect(subject.Get([]byte("key-00000333"))).To(Equal([]byte("val-00000333")))
	})
})
