// darwind3 handles the real time push port feed
package darwind3

type DarwinD3 struct {

}

// Process a messafe
type Processor interface {
  Process( *DarwinD3, *Pport ) error
}
