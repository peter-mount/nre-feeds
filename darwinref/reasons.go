package darwinref

// A reason, shared by LateRunningReasons and CancellationReasons
type Reason struct {
  Code        int               `xml:"code,attr"`
  Text        string            `xml:"reasontext,attr"`
}

// Return the cancellation reason
func (r *DarwinReference) getCancellationReason( i int ) string {
  if val, ok := r.cancellationReasons[ i ]; ok {
    return val
  }
  return ""
}

// Return the late reason
func (r *DarwinReference) getLateReason( i int ) string {
  if val, ok := r.lateRunningReasons[ i ]; ok {
    return val
  }
  return ""
}
