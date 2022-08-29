// Copyright 2020-2022 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jetstream

import (
	"errors"
	"fmt"
)

type (
	// JetStreamError is an error result that happens when using JetStream.
	// In case of client-side error, `APIError()` returns nil
	JetStreamError interface {
		APIError() *APIError
		error
	}

	jsError struct {
		apiErr  *APIError
		message string
	}

	// APIError is included in all API responses if there was an error.
	APIError struct {
		Code        int       `json:"code"`
		ErrorCode   ErrorCode `json:"err_code"`
		Description string    `json:"description,omitempty"`
	}

	// ErrorCode represents `error_code` returned in response from JetStream API
	ErrorCode uint16
)

const (
	JSErrCodeJetStreamNotEnabledForAccount ErrorCode = 10039
	JSErrCodeJetStreamNotEnabled           ErrorCode = 10076

	JSErrCodeStreamNotFound  ErrorCode = 10059
	JSErrCodeStreamNameInUse ErrorCode = 10058

	JSErrCodeConsumerNotFound      ErrorCode = 10014
	JSErrCodeConsumerNameExists    ErrorCode = 10013
	JSErrCodeConsumerAlreadyExists ErrorCode = 10105

	JSErrCodeMessageNotFound ErrorCode = 10037

	JSErrCodeBadRequest ErrorCode = 10003
)

var (
	// API errors

	// ErrJetStreamNotEnabled is an error returned when JetStream is not enabled for an account.
	ErrJetStreamNotEnabled JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeJetStreamNotEnabled, Description: "jetstream not enabled", Code: 503}}

	// ErrJetStreamNotEnabledForAccount is an error returned when JetStream is not enabled for an account.
	ErrJetStreamNotEnabledForAccount JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeJetStreamNotEnabledForAccount, Description: "jetstream not enabled for account", Code: 503}}

	// ErrStreamNotFound is an error returned when stream with given name does not exist.
	ErrStreamNotFound JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeStreamNotFound, Description: "stream not found", Code: 404}}

	// ErrStreamNameAlreadyInUse is returned when a stream with given name already exists and has a different configuration
	ErrStreamNameAlreadyInUse JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeStreamNameInUse, Description: "stream name already in use", Code: 400}}

	// ErrConsumerNotFound is an error returned when consumer with given name does not exist.
	ErrConsumerNotFound JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeConsumerNotFound, Description: "consumer not found", Code: 404}}

	// ErrMsgNotFound is returned when message with provided sequence number does npt exist.
	ErrMsgNotFound JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeMessageNotFound, Description: "message not found", Code: 404}}

	// ErrBadRequest is returned when invalid request is sent to JetStream API.
	ErrBadRequest JetStreamError = &jsError{apiErr: &APIError{ErrorCode: JSErrCodeBadRequest, Description: "bad request", Code: 400}}

	// Client errors

	// ErrConsumerNotFound is an error returned when consumer with given name does not exist.
	ErrConsumerNameAlreadyInUse JetStreamError = &jsError{message: "consumer name already in use"}

	// ErrInvalidJSAck is returned when JetStream ack from message publish is invalid.
	ErrInvalidJSAck JetStreamError = &jsError{message: "invalid jetstream publish response"}

	// ErrStreamNameRequired is returned when the provided stream name is empty.
	ErrStreamNameRequired JetStreamError = &jsError{message: "stream name is required"}

	// ErrConsumerNameRequired is returned when the provided consumer durable name is empty,
	ErrConsumerNameRequired JetStreamError = &jsError{message: "consumer name is required"}

	// ErrMsgAlreadyAckd is returned when attempting to acknowledge message more than once.
	ErrMsgAlreadyAckd JetStreamError = &jsError{message: "message was already acknowledged"}

	// ErrNoStreamResponse is returned when there is no response from stream (e.g. no responders error).
	ErrNoStreamResponse JetStreamError = &jsError{message: "no response from stream"}

	// ErrNotJSMessage is returned when attempting to get metadata from non JetStream message .
	ErrNotJSMessage JetStreamError = &jsError{message: "not a jetstream message"}

	// ErrInvalidStreamName is returned when the provided stream name is invalid (contains '.').
	ErrInvalidStreamName JetStreamError = &jsError{message: "invalid stream name"}

	// ErrInvalidConsumerName is returned when the provided consumer name is invalid (contains '.').
	ErrInvalidConsumerName JetStreamError = &jsError{message: "invalid consumer name"}

	// ErrNoMessages is returned when no messages are currectly available for a consumer
	ErrNoMessages = &jsError{message: "no messages"}

	// ErrHandlerRequired is returned when no handler func is provided in Stream()
	ErrHandlerRequired = &jsError{message: "handler cannot be empty"}

	// ErrEndOfData is returned when iterating over paged API from JetStream reaches end of data
	ErrEndOfData = errors.New("nats: end of data reached")

	// ErrNoHeartbeat is received when no message is received in IdleHeartbeat time (if set)
	ErrNoHeartbeat = &jsError{message: "no heartbeat received, canceling subscription"}

	// ErrConsumerHasActiveSubscription is returned when a consumer is already subscribed to a stream
	ErrConsumerHasActiveSubscription = &jsError{message: "consumer has active subscription"}

	// ErrMsgNotBound is returned when given message is not bound to any subscription
	ErrMsgNotBound = &jsError{message: "message is not bound to subscription/connection"}

	// ErrMsgNoReply is returned when attempting to reply to a message without a reply subject
	ErrMsgNoReply = &jsError{message: "message does not have a reply"}

	// ErrMsgDeleteUnsuccessful is returned when an attempt to delete a message is unsuccessful
	ErrMsgDeleteUnsuccessful = &jsError{message: "message deletion unsuccessful"}

	// ErrAsyncPublishReplySubjectSet is returned when reply subject is set on async message publish
	ErrAsyncPublishReplySubjectSet = &jsError{message: "reply subject should be empty"}

	// ErrTooManyStalledMsgs is returned when too many outstanding async messages are waiting for ack
	ErrTooManyStalledMsgs = &jsError{message: "stalled with too many outstanding async published messages"}
)

// Error prints the JetStream API error code and description
func (e *APIError) Error() string {
	return fmt.Sprintf("nats: API error %d: %s", e.ErrorCode, e.Description)
}

// APIError implements the JetStreamError interface.
func (e *APIError) APIError() *APIError {
	return e
}

// Is matches against an APIError.
func (e *APIError) Is(err error) bool {
	if e == nil {
		return false
	}
	// Extract internal APIError to match against.
	var aerr *APIError
	ok := errors.As(err, &aerr)
	if !ok {
		return ok
	}
	return e.ErrorCode == aerr.ErrorCode
}

func (err *jsError) APIError() *APIError {
	return err.apiErr
}

func (err *jsError) Error() string {
	if err.apiErr != nil && err.apiErr.Description != "" {
		return err.apiErr.Error()
	}
	return fmt.Sprintf("nats: %s", err.message)
}

func (err *jsError) Unwrap() error {
	// Allow matching to embedded APIError in case there is one.
	if err.apiErr == nil {
		return nil
	}
	return err.apiErr
}
