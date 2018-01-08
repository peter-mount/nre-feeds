// Debug output of a Journey
package darwintimetable

import (
  "bytes"
  //"log"
  "strconv"
  "fmt"
)

func field( b *bytes.Buffer, n string, s string ) {
  for i := len( n ); i < 15; i++ {
    b.WriteString( " " )
  }
  b.WriteString( n )
  b.WriteString( " " )
  b.WriteString( s )
  b.WriteString( "\n" )
  //log.Println( n, s )
}

func fieldi( b *bytes.Buffer, n string, s int) {
  field( b, n, strconv.FormatInt( int64( s ), 10 ) )
}

func fieldb( b *bytes.Buffer, n string, s bool) {
  if s {
    field( b, n, "true" )
  } else {
    field( b, n, "false" )
  }
}

func pad( b *bytes.Buffer, s string, l int ) {
  b.WriteString( s )
  for i := len( s ); i < l; i++ {
    b.WriteString( " " )
  }
  b.WriteString( " | " )
}

func padx( b *bytes.Buffer, l int ) {
  pad(b,"",l)
}

func padb( b *bytes.Buffer, s bool, l int ) {
  if s {
    pad( b, "t", l )
  } else {
    pad( b, "f", l )
  }
}

func padi( b *bytes.Buffer, s int, l int ) {
  pad( b, strconv.FormatInt( int64( s ), 10 ), l )
}

func (j *Journey) String() string {
  var b bytes.Buffer

  field( &b, "RID", j.RID )
  field( &b, "UID", j.UID )
  field( &b, "TrainID", j.TrainID )
  field( &b, "SSD", j.SSD )
  field( &b, "Toc", j.Toc )
  field( &b, "TrainCat", j.TrainCat )
  fieldb( &b, "Passenger", j.Passenger )
  fieldi( &b, "Cancel Reason", j.CancelReason )

  b.WriteString( "| " )
  pad( &b, "Tiploc", 8 )
  pad( &b, "Act", 4 )
  pad( &b, "PAct", 4 )
  pad( &b, "Can", 3 )
  pad( &b, "Plat", 4 )
  pad( &b, "PTA", 5 ) // pta
  pad( &b, "PTD", 5 ) // ptd
  pad( &b, "WTA", 8 )
  pad( &b, "WTD", 8 )
  pad( &b, "WTP", 8 ) // wtp
  pad( &b, "FDest", 8 ) // FalseDest
  pad( &b, "RDel", 5 ) // RDelay
  b.WriteString( "\n" )

  for _, s := range j.Schedule {
    b.WriteString( fmt.Sprintf( "%v\n", s ) )
  }

  return b.String()
}

func (l *OPOR) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  padx( &b, 5 ) // pta
  padx( &b, 5 ) // ptd
  pad( &b, l.Wta, 8 )
  pad( &b, l.Wtd, 8 )
  padx( &b, 8 ) // wtp
  padx( &b, 8 ) // FalseDest
  padx( &b, 5 ) // RDelay

  return b.String()
}

func (l *OR) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  pad( &b, l.Pta, 5 )
  pad( &b, l.Ptd, 5 )
  pad( &b, l.Wta, 8 )
  pad( &b, l.Wtd, 8 )
  padx( &b, 8 ) // wtp
  pad( &b, l.FalseDest, 8 )
  padx( &b, 5 ) // RDelay

  return b.String()
}

func (l *IP) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  pad( &b, l.Pta, 5 )
  pad( &b, l.Ptd, 5 )
  pad( &b, l.Wta, 8 )
  pad( &b, l.Wtd, 8 )
  padx( &b, 8 ) // wtp
  pad( &b, l.FalseDest, 8 )
  pad( &b, l.RDelay, 5 )

  return b.String()
}

func (l *OPIP) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  padx( &b, 5 ) // pta
  padx( &b, 5 ) // ptd
  pad( &b, l.Wta, 8 )
  pad( &b, l.Wtd, 8 )
  padx( &b, 8 ) // wtp
  padx( &b, 8 ) // FalseDest
  pad( &b, l.RDelay, 5 )

  return b.String()
}

func (l *PP) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  padx( &b, 5 ) // pta
  padx( &b, 5 ) // ptd
  padx( &b, 8 ) // wta
  padx( &b, 8 ) // wtd
  pad( &b, l.Wtp, 8 )
  padx( &b, 8 ) // FalseDest
  pad( &b, l.RDelay, 5 )

  return b.String()
}

func (l *DT) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  pad( &b, l.Pta, 5 )
  pad( &b, l.Ptd, 5 )
  pad( &b, l.Wta, 8 )
  pad( &b, l.Wtd, 8 )
  padx( &b, 8 ) // wtp
  padx( &b, 8 ) // FalseDest
  pad( &b, l.RDelay, 5 )

  return b.String()
}

func (l *OPDT) String() string {
  var b bytes.Buffer
  b.WriteString( "| " )
  pad( &b, l.Tiploc, 8 )
  pad( &b, l.Act, 4 )
  pad( &b, l.PlanAct, 4 )
  padb( &b, l.Cancelled, 3 )
  pad( &b, l.Platform, 4 )
  padx( &b, 5 ) // pta
  padx( &b, 5 ) // ptd
  pad( &b, l.Wta, 8 )
  pad( &b, l.Wtd, 8 )
  padx( &b, 8 ) // wtp
  padx( &b, 8 ) // FalseDest
  pad( &b, l.RDelay, 5 )

  return b.String()
}
