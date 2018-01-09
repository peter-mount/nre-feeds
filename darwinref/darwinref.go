package darwinref

import (
  "fmt"
)

// Useful for debugging
func (ref *DarwinReference) String() string {
  return fmt.Sprintf(
    "DarwinReference[TimetableId=%s, Tiploc=%d, Crs=%d, Toc=%d, LateRunningReasons=%d, CancellationReasons=%d, CISSource=%d, Via=%d]",
    ref.timetableId,
    len( ref.Tiploc ),
    len( ref.Crs ),
    len( ref.Toc ),
    len( ref.LateRunningReasons ),
    len( ref.CancellationReasons ),
    len( ref.CISSource ),
    len( ref.via ) )
}

// Return's the timetableId for this reference dataset
func (ref *DarwinReference) TimetableId() string {
  return ref.timetableId
}
