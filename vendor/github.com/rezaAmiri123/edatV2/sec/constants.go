package sec

import "github.com/rezaAmiri123/edatV2/am"

const (
	notCompensating = false
	isCompensating  = true
)

const (
	SagaCommandIDHdr   = am.CommandHdrPrefix + "SAGA_ID"
	SagaCommandNameHdr = am.CommandHdrPrefix + "SAGA_NAME"

	SagaReplyIDHdr   = am.ReplyHdrPrefix + "SAGA_ID"
	SagaReplyNameHdr = am.ReplyHdrPrefix + "SAGA_NAME"
)
