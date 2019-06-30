package cdb64_test

import (
	"github.com/bsm/cdb64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	var subject *cdb64.Reader

	get := func(key []byte) (string, error) {
		val, err := subject.Get(key)
		return string(val), err
	}

	BeforeEach(func() {
		w, err := cdb64.Create(testDir + "/test.cdb")
		Expect(err).NotTo(HaveOccurred())
		defer w.Close()

		Expect(seedData(w, 1000)).To(Succeed())

		// seed some more exotic entries
		Expect(w.Put(nil, []byte("blank"))).To(Succeed())
		Expect(w.Put([]byte("blank"), nil)).To(Succeed())
		Expect(w.Put([]byte("key-00000333"), []byte("duplicate"))).To(Succeed())

		subject, err = w.Freeze()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(subject.Close()).To(Succeed())
	})

	It("should Get", func() {
		Expect(subject.Get(nil)).To(Equal([]byte("blank")))
		Expect(subject.Get([]byte{})).To(Equal([]byte("blank")))
		Expect(subject.Get([]byte("missing"))).To(BeNil())
		Expect(subject.Get([]byte("blank"))).To(Equal([]byte{}))
		Expect(get([]byte("key-00000005"))).To(Equal("val-00000005"))
		Expect(get([]byte("key-00000333"))).To(Equal("val-00000333"))
	})

	It("should iterate", func() {
		iter := subject.Iterator()

		var keys, vals []string
		for iter.Next() {
			keys = append(keys, string(iter.Key()))
			vals = append(vals, string(iter.Value()))
		}
		Expect(len(keys)).To(Equal(1003))
		Expect(len(vals)).To(Equal(1003))

		Expect(keys[:3]).To(Equal([]string{"key-00000001", "key-00000003", "key-00000005"}))
		Expect(vals[:3]).To(Equal([]string{"val-00000001", "val-00000003", "val-00000005"}))
		Expect(keys[165:168]).To(Equal([]string{"key-00000331", "key-00000333", "key-00000335"}))
		Expect(vals[165:168]).To(Equal([]string{"val-00000331", "val-00000333", "val-00000335"}))
		Expect(keys[1000:]).To(Equal([]string{"", "blank", "key-00000333"}))
		Expect(vals[1000:]).To(Equal([]string{"blank", "", "duplicate"}))
	})
})
