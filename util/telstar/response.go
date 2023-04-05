package telstar

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Response is returned to Telstar
type Response struct {
	pageNumber int
	frameId    byte
	frames     []*Frame
	dynamic    bool
}

func NewResponse() *Response {
	return &Response{frameId: 'a'}
}

// PageNumber sets the page number of the response
func (r *Response) PageNumber(n int) *Response {
	r.pageNumber = n
	return r
}

// FrameId sets the initial frame id. Default is "a" but can be "b"
func (r *Response) FrameId(n byte) *Response {
	if n < 'a' || n > 'z' {
		panic(fmt.Errorf("invalid frameId %c", n))
	}
	r.frameId = n
	return r
}

// Dynamic sets the response as dynamic - for use with our
// fork of Telstar only
func (r *Response) Dynamic() *Response {
	r.dynamic = true
	return r
}

// AddFrame adds a frame to the response
func (r *Response) addFrame(f *Frame) *Response {
	r.frames = append(r.frames, f)
	return r
}

func (r *Response) Build() (string, error) {
	if len(r.frames) == 0 {
		return "", errors.New("no frames in response")
	}

	frameId := r.frameId
	for fn, f := range r.frames {
		// Ensure frames have the correct PID's
		f.PID.PageNumber = r.pageNumber
		f.PID.FrameId = string(frameId)
		frameId++
		if frameId > 'z' {
			return "", errors.New("too many frames for response")
		}

		if fn == 0 && r.dynamic {
			// Only set the first frame as dynamic if we are in that mode
			f.FrameType = "dynamic"
			// Ensure we have this set pointing back to us so subsequent calls to the page refresh
			f.ResponseData = ResponseData{
				Fields: nil,
				Action: ResponseAction{
					Exec:            os.Args[0],
					Args:            os.Args[1:],
					PostActionFrame: f.PID,
				},
			}
		} else if f.FrameType == "" {
			// If not set then set it to information
			f.FrameType = "information"
		}

		// Ensure we have NavMessage set, if not then it breaks Telstar
		if f.NavMessage == "" {
			f.NavMessage = "[B][n][Y]Select item or[W]*page# : [_+]"
		}
		if f.NavMessageNotFound == "" {
			f.NavMessageNotFound = "[B][n][Y]Page not Found :[W]"
		}

		// Ensure we have a routing table TODO set this to something else?
		if len(f.RoutingTable) == 0 {
			f.RoutingTable = []int{r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber, r.pageNumber}
		}

		// Force these
		f.AuthorIdOld = "nre-feeds"
		f.Visible = true
	}

	var output strings.Builder
	for _, f := range r.frames {
		outputText, err := f.Dump()
		if err != nil {
			return "", err
		}

		_, err = output.Write(outputText)
		if err != nil {
			return "", err
		}

		_, err = output.WriteString(",")
		if err != nil {
			return "", err
		}
	}

	outs := output.String()
	return fmt.Sprintf("[%s]", outs[:len(outs)-1]), nil
}
