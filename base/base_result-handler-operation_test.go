// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	kvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
)

var _ = Describe("ResultHandlerTxOperation", func() {
	var ctx context.Context
	var err error
	var resultHandlerTxOperation base.ResultHandlerTx
	var baseResultHandler *libmocks.BaseResultHandlerTx
	var result base.Result
	BeforeEach(func() {
		ctx = context.Background()
		result = base.Result{
			Operation: "my-operation",
		}
		baseResultHandler = &libmocks.BaseResultHandlerTx{}
		resultHandlerTxOperation = base.ResultHandlerTxOperation{
			"my-operation": baseResultHandler,
		}
	})
	Context("HandleResult", func() {
		JustBeforeEach(func() {
			err = resultHandlerTxOperation.HandleResult(ctx, &kvmocks.Tx{}, result)
		})
		Context("default operation", func() {
			BeforeEach(func() {
				result.Operation = "other-operation"
			})
			It("does not call resultHandler", func() {
				Expect(baseResultHandler.HandleResultCallCount()).To(Equal(0))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("existing operation", func() {
			BeforeEach(func() {
			})
			It("does call resultHandler", func() {
				Expect(baseResultHandler.HandleResultCallCount()).To(Equal(1))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
