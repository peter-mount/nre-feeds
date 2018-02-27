package darwinref

import (
  "fmt"
)

// GetStations returns all Location's with a CRS code
func (c *DarwinRefClient) getReason( api string , reason int ) ( *Reason, error ) {
  if reason <= 0 {
    return nil, nil
  }

  var res *Reason

  if found, err := c.get( fmt.Sprintf( api, reason ), &res ); err != nil {
    return nil, err
  } else if found {
    return res, nil
  } else {
    return nil, nil
  }
}

func (c *DarwinRefClient) GetCancelledReason( reason int ) ( *Reason, error ) {
  return c.getReason( "/reason/cancelled/%d", reason )
}

func (c *DarwinRefClient) GetLateReason( reason int ) ( *Reason, error ) {
  return c.getReason( "/reason/late/%d", reason )
}
