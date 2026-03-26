// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	stderrors "errors"

	kvmocks "github.com/bborbe/kv/mocks"
	"github.com/bborbe/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("CommandObjectExecutorResultSender", func() {
	var ctx context.Context
	var err error
	var commandObjectExecutorResultSender raw.CommandObjectExecutor
	var commandObjectExecutor *libmocks.RawCommandObjectExecutor
	var resultObjectSender *libmocks.RawResultObjectSender
	var commandObject raw.CommandObject
	BeforeEach(func() {
		ctx = context.Background()
		commandObject = raw.CommandObject{}
		commandObjectExecutor = &libmocks.RawCommandObjectExecutor{}
		commandObjectExecutor.CommandOperationReturns("my-command")
		commandObjectExecutor.SendResultEnabledReturns(true)
		commandObjectExecutor.HandleCommandReturns(
			base.EventID("1337").Ptr(),
			base.Event{"foo": "bar"},
			nil,
		)
		resultObjectSender = &libmocks.RawResultObjectSender{}
	})
	JustBeforeEach(func() {
		commandObjectExecutorResultSender = raw.NewCommandObjectExecutorResultSender(
			commandObjectExecutor,
			resultObjectSender,
			log.DefaultSamplerFactory,
		)
	})
	It("returns correct CommandOperation", func() {
		Expect(
			commandObjectExecutorResultSender.CommandOperation(),
		).To(Equal(base.CommandOperation("my-command")))
	})
	It("returns correct SendResultEnabled", func() {
		Expect(commandObjectExecutorResultSender.SendResultEnabled()).To(BeTrue())
	})
	Context("CommandOperation", func() {
		var id *base.EventID
		var event base.Event
		JustBeforeEach(func() {
			id, event, err = commandObjectExecutorResultSender.HandleCommand(
				ctx,
				&kvmocks.Tx{},
				commandObject,
			)
		})
		Context("success send no result", func() {
			BeforeEach(func() {
				commandObjectExecutor.SendResultEnabledReturns(false)
			})
			It("returns correct id", func() {
				Expect(id).NotTo(BeNil())
				Expect(*id).To(Equal(base.EventID("1337")))
			})
			It("returns correct event", func() {
				Expect(event).NotTo(BeNil())
				Expect(event).To(Equal(base.Event{"foo": "bar"}))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("send no result", func() {
				Expect(resultObjectSender.SendCallCount()).To(Equal(0))
			})
		})
		Context("success send result", func() {
			BeforeEach(func() {
				commandObjectExecutor.SendResultEnabledReturns(true)
			})
			It("returns correct id", func() {
				Expect(id).NotTo(BeNil())
				Expect(*id).To(Equal(base.EventID("1337")))
			})
			It("returns correct event", func() {
				Expect(event).NotTo(BeNil())
				Expect(event).To(Equal(base.Event{"foo": "bar"}))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("send result", func() {
				Expect(resultObjectSender.SendCallCount()).To(Equal(1))
			})
		})
		Context("failure send result", func() {
			BeforeEach(func() {
				commandObjectExecutor.SendResultEnabledReturns(false)
				commandObjectExecutor.HandleCommandReturns(nil, nil, stderrors.New("banana"))
			})
			It("returns correct id", func() {
				Expect(id).To(BeNil())
			})
			It("returns correct event", func() {
				Expect(event).To(BeNil())
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("send result", func() {
				Expect(resultObjectSender.SendCallCount()).To(Equal(1))
			})
		})
		Context("skipped send no result", func() {
			BeforeEach(func() {
				commandObjectExecutor.SendResultEnabledReturns(false)
				commandObjectExecutor.HandleCommandReturns(nil, nil, raw.CommandObjectSkippedError)
			})
			It("returns correct id", func() {
				Expect(id).To(BeNil())
			})
			It("returns correct event", func() {
				Expect(event).To(BeNil())
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("send result", func() {
				Expect(resultObjectSender.SendCallCount()).To(Equal(0))
			})
		})
		Context("skipped send no result", func() {
			BeforeEach(func() {
				commandObjectExecutor.SendResultEnabledReturns(true)
				commandObjectExecutor.HandleCommandReturns(nil, nil, raw.CommandObjectSkippedError)
			})
			It("returns correct id", func() {
				Expect(id).To(BeNil())
			})
			It("returns correct event", func() {
				Expect(event).To(BeNil())
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("send result", func() {
				Expect(resultObjectSender.SendCallCount()).To(Equal(0))
			})
		})
	})
})
