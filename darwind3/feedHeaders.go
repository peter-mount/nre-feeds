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
	h.SequenceNumber = h.intHeader(msg, "SequenceNumber")
	h.MessageType = h.stringHeader(msg, "MessageType")
	h.PushPortSequence = h.stringHeader(msg, "PushPortSequence")
}

func (h *FeedHeaders) intHeader(msg amqp.Delivery, s string) int32 {
	sn := msg.Headers[s]
	if sn != nil {
		return sn.(int32)
	}
	return 0
}

func (h *FeedHeaders) stringHeader(msg amqp.Delivery, s string) string {
	sn := msg.Headers[s]
	if sn != nil {
		return sn.(string)
	}
	return ""
}
