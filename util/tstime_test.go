package util

import (
	"encoding/json"
	"testing"
)

const (
	ts_et    = "01:01:01"
	ts_etmin = "02:02:02"
	ts_wet   = "03:03:03"
	ts_at    = "04:04:04"
)

func TestTSTime_Equals(t *testing.T) {
	mk := func() *TSTime {
		a := &TSTime{}
		a.ET = NewWorkingTime(ts_et)
		a.ETMin = NewWorkingTime(ts_etmin)
		a.WET = NewWorkingTime(ts_wet)
		a.AT = NewWorkingTime(ts_at)
		return a
	}

	a := mk()
	b := mk()

	// Test times are ok
	if !a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}

	// Now test each other parameter one by one. These should not equal

	b.ETUnknown = true
	if a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}
	a.ETUnknown = b.ETUnknown

	b.ATRemoved = true
	if a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}
	a.ATRemoved = b.ATRemoved

	b.Delayed = true
	if a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}
	a.Delayed = true

	b.Src = "src"
	if a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}
	a.Src = b.Src

	b.SrcInst = "Inst"
	if a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}
	a.SrcInst = b.SrcInst

	// Final equality check now both should be the same
	if !a.Equals(b) {
		t.Errorf("expected %v, got %v", a, b)
	}
}

// Test we Marshal/Unmarshal JSON correctly
func TestTSTime_JSON(t *testing.T) {

	gett := func(s string) *WorkingTime {
		if s == "00:00:00" {
			return nil
		}
		return NewWorkingTime(s)
	}

	tst := func(f func(v *TSTime, s string)) bool {
		return runHHMMSS_TimeSeries(t, func(s string) bool {
			v := &TSTime{}
			f(v, s)
			return tstime_test(t, v)
		})
	}

	tst(func(v *TSTime, s string) {
		v.ET = gett(s)
	})
	tst(func(v *TSTime, s string) {
		v.ETMin = gett(s)
	})
	tst(func(v *TSTime, s string) {
		v.WET = gett(s)
	})
	tst(func(v *TSTime, s string) {
		v.AT = gett(s)
	})
}

func tstime_testTime(t *testing.T, s string, a *WorkingTime, b *WorkingTime) bool {
	if a == nil {
		if b != nil {
			t.Errorf("%s expected nil, got %v", s, b)
			return true
		}
	} else if a.IsZero() != b.IsZero() {
		t.Errorf("%s expected zero %v, got %v", s, a.IsZero(), b.IsZero())
		return true
	} else if !a.Equals(b) {
		t.Errorf("%s expected %v, got %v", s, a, b)
		return true
	}
	return false
}

func tstime_test(t *testing.T, a *TSTime) bool {

	if b, err := json.Marshal(a); err != nil {
		t.Errorf("%v failed to marshal: %v", a, err)
		return true
	} else {

		// unmarshal back
		c := &TSTime{}
		if err := json.Unmarshal(b, c); err != nil {
			t.Errorf("%v failed to marshal: %v", a, err)
			return true
		}

		return tstime_testTime(t, "ET", a.ET, c.ET) ||
			tstime_testTime(t, "ETMin", a.ETMin, c.ETMin) ||
			tstime_testTime(t, "WET", a.WET, c.WET) ||
			tstime_testTime(t, "AT", a.AT, c.AT)
	}

}
