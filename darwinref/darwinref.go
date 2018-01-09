package darwinref

import (
  "fmt"
)

// Useful for debugging
func (ref *DarwinReference) String() string {
  return fmt.Sprintf(
    "DarwinReference[TimetableId=%s, Tiploc=%d, Crs=%d, Toc=%d, LateRunningReasons=%d, CancellationReasons=%d, CISSource=%d, Via=%d]",
    ref.timetableId,
    len( ref.tiploc ),
    len( ref.crs ),
    len( ref.toc ),
    len( ref.lateRunningReasons ),
    len( ref.cancellationReasons ),
    len( ref.cisSource ),
    len( ref.via ) )
}

// Return's the timetableId for this reference dataset
func (ref *DarwinReference) TimetableId() string {
  return ref.timetableId
}
