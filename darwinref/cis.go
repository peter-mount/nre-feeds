package darwinref

// The CIS Source
func (r *DarwinReference) getCISSource( s string ) string {
  if val, ok := r.cisSource[ s ]; ok {
    return val
  }
  return ""
}
