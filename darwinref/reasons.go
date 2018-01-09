package darwinref

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
