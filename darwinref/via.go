package darwinref

// Are two Via's equal
func (v *Via) Equals( o *Via ) bool {
  if o == nil {
    return false
  }
  return v.At == o.At && v.Dest == o.Dest && v.Loc1 == o.Loc1 && v.Loc2 == o.Loc2
}

func (v *Via) String() string {
  return "Via[At=" + v.At +", Dest=" + v.Dest +", Loc1=" + v.Loc1 +", Loc2=" + v.Loc2 +", Text=" + v.Text + "]"
}

// Return all Via's at a location with a specific destination
// at   CRS code of location
// dest Tiploc of train's destination
func (r *DarwinReference) GetViaAt( at string, dest string ) ( []*Via, bool ) {
  val, ok := r.via[ at + "," + dest ]
  return val, ok
}
