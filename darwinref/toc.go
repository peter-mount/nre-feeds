package darwinref

// Lookup a Toc from it's ATOC code
func (r *DarwinReference) getToc( t string ) ( *Toc, bool ) {
  val, ok := r.toc[ t ]
  return val, ok
}
