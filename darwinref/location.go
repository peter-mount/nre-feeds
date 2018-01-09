package darwinref

// Return a *Location for a tiploc
func (r *DarwinReference) getTiploc( t string ) ( *Location, bool ) {
  val, ok := r.tiploc[ t ]
  return val, ok
}

// Return a []*Location for a CRS / 3Alpha code
func (r *DarwinReference) getCrs( c string ) ( []*Location, bool ) {
  val, ok := r.crs[ c ]
  return val, ok
}
