// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"strconv"
	"sync"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("IdentifierGenerator", func() {
	Context("IdentifierGeneratorFunc", func() {
		It("implements IdentifierGenerator interface", func() {
			generator := base.IdentifierGeneratorFunc[string](func() string {
				return "test-id"
			})

			result := generator.NewIdentifier()
			Expect(result).To(Equal("test-id"))
		})

		It("calls the underlying function", func() {
			callCount := 0
			generator := base.IdentifierGeneratorFunc[string](func() string {
				callCount++
				return "test-id-" + strconv.Itoa(callCount)
			})

			result1 := generator.NewIdentifier()
			result2 := generator.NewIdentifier()

			Expect(result1).To(Equal("test-id-1"))
			Expect(result2).To(Equal("test-id-2"))
			Expect(callCount).To(Equal(2))
		})
	})

	Context("NewIdentifierGeneratorUUID", func() {
		Context("string type", func() {
			var generator base.IdentifierGenerator[string]

			BeforeEach(func() {
				generator = base.NewIdentifierGeneratorUUID[string]()
			})

			It("generates valid UUID strings", func() {
				id1 := generator.NewIdentifier()
				id2 := generator.NewIdentifier()

				Expect(id1).ToNot(BeEmpty())
				Expect(id2).ToNot(BeEmpty())
				Expect(id1).ToNot(Equal(id2))

				// Verify they are valid UUIDs
				_, err1 := uuid.Parse(id1)
				_, err2 := uuid.Parse(id2)
				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
			})

			It("generates unique IDs", func() {
				ids := make(map[string]bool)
				for i := 0; i < 100; i++ {
					id := generator.NewIdentifier()
					Expect(ids[id]).To(BeFalse(), "ID should be unique: %s", id)
					ids[id] = true
				}
			})
		})

		Context("byte slice type", func() {
			var generator base.IdentifierGenerator[[]byte]

			BeforeEach(func() {
				generator = base.NewIdentifierGeneratorUUID[[]byte]()
			})

			It("generates valid UUID byte slices", func() {
				id1 := generator.NewIdentifier()
				id2 := generator.NewIdentifier()

				Expect(id1).ToNot(BeEmpty())
				Expect(id2).ToNot(BeEmpty())
				Expect(id1).ToNot(Equal(id2))

				// Verify they are valid UUIDs
				_, err1 := uuid.Parse(string(id1))
				_, err2 := uuid.Parse(string(id2))
				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
			})
		})
	})

	Context("NewIdentifierGeneratorCounter", func() {
		Context("string type", func() {
			var generator base.IdentifierGenerator[string]

			BeforeEach(func() {
				generator = base.NewIdentifierGeneratorCounter[string]()
			})

			It("generates sequential counter IDs", func() {
				id1 := generator.NewIdentifier()
				id2 := generator.NewIdentifier()
				id3 := generator.NewIdentifier()

				Expect(id1).To(Equal("1"))
				Expect(id2).To(Equal("2"))
				Expect(id3).To(Equal("3"))
			})

			It("is thread-safe", func() {
				const goroutines = 10
				const idsPerGoroutine = 100

				var wg sync.WaitGroup
				idsChan := make(chan string, goroutines*idsPerGoroutine)

				for i := 0; i < goroutines; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						for j := 0; j < idsPerGoroutine; j++ {
							idsChan <- generator.NewIdentifier()
						}
					}()
				}

				wg.Wait()
				close(idsChan)

				// Collect all IDs
				ids := make(map[string]bool)
				var allIDs []string
				for id := range idsChan {
					Expect(ids[id]).To(BeFalse(), "ID should be unique: %s", id)
					ids[id] = true
					allIDs = append(allIDs, id)
				}

				Expect(len(allIDs)).To(Equal(goroutines * idsPerGoroutine))
				Expect(len(ids)).To(Equal(goroutines * idsPerGoroutine))
			})
		})

		Context("byte slice type", func() {
			var generator base.IdentifierGenerator[[]byte]

			BeforeEach(func() {
				generator = base.NewIdentifierGeneratorCounter[[]byte]()
			})

			It("generates sequential counter IDs as bytes", func() {
				id1 := generator.NewIdentifier()
				id2 := generator.NewIdentifier()
				id3 := generator.NewIdentifier()

				Expect(string(id1)).To(Equal("1"))
				Expect(string(id2)).To(Equal("2"))
				Expect(string(id3)).To(Equal("3"))
			})
		})
	})
})
