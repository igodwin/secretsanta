package utils_test

import (
	. "github.com/igodwin/secretsanta/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	const (
		expectedFirst  = "first"
		expectedSecond = "second"
		expectedThird  = "third"
	)

	var queue *Queue

	BeforeEach(func() {
		queue = &Queue{}
	})

	Context("Enqueue and Dequeue", func() {
		It("enqueues elements that may be dequeued in expected order", func() {
			queue.Enqueue(expectedFirst)
			queue.Enqueue(expectedSecond)
			queue.Enqueue(expectedThird)

			Expect(queue.Dequeue()).To(Equal(expectedFirst))
			Expect(queue.Dequeue()).To(Equal(expectedSecond))
			Expect(queue.Dequeue()).To(Equal(expectedThird))
			Expect(queue.IsEmpty()).To(BeTrue())
		})

		It("dequeue returns nil on empty queue", func() {
			Expect(queue.Dequeue()).To(BeNil())
		})
	})

	Context("Size", func() {
		It("returns expected size", func() {
			Expect(queue.Size()).To(BeNumerically("==", 0))
			queue.Enqueue(expectedFirst)
			Expect(queue.Size()).To(BeNumerically("==", 1))
			queue.Enqueue(expectedSecond)
			Expect(queue.Size()).To(BeNumerically("==", 2))
			_ = queue.Dequeue()
			Expect(queue.Size()).To(BeNumerically("==", 1))
			_ = queue.Dequeue()
			Expect(queue.Size()).To(BeNumerically("==", 0))
			Expect(queue.IsEmpty()).To(BeTrue())
		})
	})

	Context("IsEmpty", func() {
		It("returns true initially", func() {
			Expect(queue.IsEmpty()).To(BeTrue())
		})

		It("returns false if one or more enqueued", func() {
			queue.Enqueue(expectedFirst)
			Expect(queue.IsEmpty()).To(BeFalse())
			queue.Enqueue(expectedSecond)
			Expect(queue.IsEmpty()).To(BeFalse())
			queue.Enqueue(expectedThird)
			Expect(queue.IsEmpty()).To(BeFalse())
		})

		It("returns true after all elements have been dequeued", func() {
			queue.Enqueue(expectedFirst)
			queue.Enqueue(expectedSecond)
			queue.Enqueue(expectedThird)
			Expect(queue.IsEmpty()).To(BeFalse())
			_ = queue.Dequeue()
			_ = queue.Dequeue()
			_ = queue.Dequeue()
			Expect(queue.IsEmpty()).To(BeTrue())
		})
	})
})
