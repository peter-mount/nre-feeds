package darwind3

import (
	"github.com/streadway/amqp"
)

// FeedHeaders holds the relevant headers from the feed
// used by FeedStatus to detect any missing messages etc
type FeedHeaders struct {
	SequenceNumber   int32  `json:"sequenceNumber"`
	PushPortSequence string `json:"pushPortSequence"`
	MessageType      string `json:"messageType"`
}

// Populate the FeedHeaders from the inbound message
func (h *FeedHeaders) populate(msg amqp.Delivery) {
	h.SequenceNumber = msg.Headers["SequenceNumber"].(int32)
	h.MessageType = msg.Headers["MessageType"].(string)
	h.PushPortSequence = msg.Headers["PushPortSequence"].(string)
}
