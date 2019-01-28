package darwind3

import (
  "github.com/muesli/cache2go"
  "io/ioutil"
  "os"
  "time"
)

const (
  expiryTime time.Duration = 120 * time.Second
  persistSize int = 100
)

// Memory cache of schedules with disk persistance
type cache struct {
  // Schedule cache
  scheduleCache        *cache2go.CacheTable
  // Location of disk cache
  cacheDir              string
}

func (c *cache) initCache( cacheDir string ) error {
  c.cacheDir = cacheDir

  if err := os.MkdirAll( cacheDir, 0777 ); err != nil {
    return err
  }

  c.scheduleCache = cache2go.Cache( "schedules" )

  // If not in the cache then look to the disk
  c.scheduleCache.SetDataLoader( func(key interface{}, args ...interface{}) *cache2go.CacheItem {
    path, fn := c.getPath( key.(string) )
    if b, err := ioutil.ReadFile( path + "/" + fn ); err != nil {
      return nil
    } else {

      sched := ScheduleFromBytes( b )
      if sched == nil || sched.RID == "" {
        return nil
      }

      return cache2go.NewCacheItem( key, expiryTime, sched )
    }
  } )

  return nil
}

// Retrieve a schedule by it's rid
func (d *DarwinD3) GetSchedule( rid string ) *Schedule {
  val, err := d.cache.scheduleCache.Value( rid )
  if err != nil {
    sched := d.resolveSchedule( rid )
    if sched != nil {
      d.putSchedule( sched )
    }
    return sched
  }
  if err == nil {
    return val.Data().(*Schedule)
  }
  return nil
}

// Store a schedule by it's rid
func (d *DarwinD3) putSchedule( sched *Schedule ) {
  d.cache.scheduleCache.Add( sched.RID, expiryTime, sched )
  d.cache.persistSchedule( sched )
}

func (c *cache) persistSchedule( sched *Schedule ) {
  if b, err := sched.Bytes(); err == nil {
    dir, fn := c.getPath( sched.RID )
    os.MkdirAll( dir, 0777 )
    ioutil.WriteFile( dir + "/" + fn, b, 0655 )
  }
}

// Delete a schedule
func (d *DarwinD3) deleteSchedule( rid string ) {
  d.cache.scheduleCache.Delete( rid )
  path, fn := d.cache.getPath( rid )
  os.Remove( path + "/" + fn )
}

func (c *cache) getPath( rid string ) ( string, string ) {
  return c.cacheDir + "/" + rid[0:6] + "/" + rid[6:8] + "/" + rid[8:10] + "/" + rid[10:12], rid[12:]
}
