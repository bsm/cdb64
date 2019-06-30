package cdb64_test

import (
	"github.com/bsm/cdb64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iterator", func() {
	var subject *cdb64.Iterator
	var reader *cdb64.Reader

	BeforeEach(func() {
		w, err := cdb64.Create(testDir + "/test.cdb")
		Expect(err).NotTo(HaveOccurred())
		Expect(seedData(w, 10)).To(Succeed())

		reader, err = w.Freeze()
		Expect(err).NotTo(HaveOccurred())

		subject = reader.Iterator()
	})

	AfterEach(func() {
		Expect(reader.Close()).To(Succeed())
	})

	It("should iterate (re-using keys and values)", func() {
		Expect(subject.Next()).To(BeTrue())
		Expect(string(subject.Key())).To(Equal("key-00000001"))
		Expect(string(subject.Value())).To(Equal("val-00000001"))

		key, val := subject.Key(), subject.Value()
		Expect(string(key)).To(Equal("key-00000001"))
		Expect(string(val)).To(Equal("val-00000001"))

		Expect(subject.Next()).To(BeTrue())
		Expect(string(subject.Key())).To(Equal("key-00000003"))
		Expect(string(subject.Value())).To(Equal("val-00000003"))
		Expect(string(key)).To(Equal("key-00000003"))
		Expect(string(val)).To(Equal("val-00000003"))
	})
})
