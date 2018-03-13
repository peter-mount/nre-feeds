package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "sort"
  "strings"
)

type SearchResult struct {
  Crs       string      `json:"code"`
  Name      string      `json:"name"`
  Label     string      `json:"label"`
  Score     float64     `json:"score"`
  Distance  float64     `json:"distance,omitempty"`
}

func newSearchResult( l *Location, score, d float64 ) *SearchResult {
  var label = l.Name + " [" + l.Crs + "]"
  if d > 0 {
    label = fmt.Sprintf( "%s [%s] %0.1fkm", l.Name, l.Crs, d )
  } else {
    label = fmt.Sprintf( "%s [%s]", l.Name, l.Crs )
  }
  return &SearchResult{
    Crs: l.Crs,
    Name: l.Name,
    Label: label,
    Score: score,
    Distance: d,
  }
}

func (dr *DarwinReference) SearchName( term string ) ([]*SearchResult, error) {
  if len( term ) < 3 {
    return nil, nil
  }

  term = strings.ToUpper( term )

  set := make( map[string]*SearchResult )

  if err := dr.View( func( tx *bolt.Tx ) error {
    crsBucket := tx.Bucket( []byte( "DarwinCrs" ) )
    tiplocBucket := tx.Bucket( []byte( "DarwinTiploc" ) )

    return crsBucket.ForEach( func( k, v []byte ) error {
      var tpls []string
      codec.NewBinaryCodecFrom( v ).ReadStringArray( &tpls )

      appendCrs := len( term ) == 3 && string(k[:]) == term

      for _, tpl := range tpls {
        if loc, exists := dr.GetTiplocBucket( tiplocBucket, tpl ); exists {
          var score float64

          if appendCrs {
            score = 1.0
          } else {
            s := strings.ToUpper( loc.Name )
            if strings.Contains( s, term ) {
              score = float64(len( term )) / float64(len( s ))
            }
          }

          if score > 0.0 && loc.IsPublic() {
            if _, exists := set[loc.Crs]; !exists {
              set[ loc.Crs ] = newSearchResult(loc,score,0.0)
            }
          }
        }
      }
      return nil
    } )
  }); err != nil {
    return nil, err
  }

  // Get slice of values then sort by score descending
  var result []*SearchResult
  for _, l := range set {
    result = append( result, l )
  }
  sort.SliceStable( result, func(i, j int) bool {
    // Compare scores at 2 decimal places
    s1 := int32(100*result[i].Score)
    s2 := int32(100*result[j].Score)
    if s1 == s2 {
      // Sort identical scores alphabetically
      return strings.ToUpper( result[i].Name ) < strings.ToUpper( result[j].Name )
    } else {
      // Sort scores by descending score order
      return s1 > s2
    }
  })

  return result, nil
}

func (dr *DarwinReference) SearchHandler( r *rest.Rest ) error {
  if results, err := dr.SearchName( r.Var( "term" ) ); err != nil {
    return err
  } else {
    r.Status( 200 ).Value( results )
    return nil
  }
}
