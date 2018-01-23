package darwinref

// A rail operator
type Toc struct {
  Toc         string            `xml:"toc,attr"`
  Name        string            `xml:"tocname,attr"`
  Url         string            `xml:"url,attr"`
}

// Lookup a Toc from it's ATOC code
func (r *DarwinReference) getToc( t string ) ( *Toc, bool ) {
  val, ok := r.toc[ t ]
  return val, ok
}
