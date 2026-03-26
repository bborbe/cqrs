// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

type MyEvent base.Event

func (m MyEvent) Validate(ctx context.Context) error {
	return nil
}

var _ = Describe("Command Creator", func() {
	var ch <-chan base.RequestID
	var commandCreator base.CommandCreator
	var object base.Event
	var command base.Command
	var ctx context.Context
	var cancel context.CancelFunc
	var id base.EventID
	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		ch = base.RequestIDChannel(ctx)
		commandCreator = base.NewCommandCreator(ch)
		object = base.Event{}
		id = "123"
	})
	AfterEach(func() {
		cancel()
	})
	Context("NewCommand", func() {
		BeforeEach(func() {
			command = commandCreator.NewCommand("my-command", "my-user", id, object)
		})
		It("returns command with operation", func() {
			Expect(command.Operation).To(Equal(base.CommandOperation("my-command")))
		})
		It("returns command with request id", func() {
			Expect(command.RequestID).NotTo(Equal(base.RequestID("")))
		})
		It("returns command with object", func() {
			Expect(command.Data).To(Equal(object))
		})
		It("has request time", func() {
			Expect(command.RequestTime).NotTo(Equal(0))
		})
	})
	Context("CreateCommand", func() {
		BeforeEach(func() {
			command = commandCreator.CreateCommand("my-user", object)
		})
		It("returns command with operation", func() {
			Expect(command.Operation).To(Equal(base.CreateOperation))
		})
	})
	Context("DeleteCommand", func() {
		BeforeEach(func() {
			command = commandCreator.DeleteCommand("my-user", id)
		})
		It("returns command with operation", func() {
			Expect(command.Operation).To(Equal(base.DeleteOperation))
		})
	})
	Context("UpdateCommand", func() {
		BeforeEach(func() {
			command = commandCreator.UpdateCommand("my-user", id, object)
		})
		It("returns command with operation", func() {
			Expect(command.Operation).To(Equal(base.UpdateOperation))
		})
	})
	Context("PatchCommand", func() {
		BeforeEach(func() {
			command = commandCreator.PatchCommand("my-user", id, object)
		})
		It("returns command with operation", func() {
			Expect(command.Operation).To(Equal(base.PatchOperation))
		})
	})
})
