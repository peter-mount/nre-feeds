package telstar

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	type args struct {
		s string
		l int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test",
			args: args{
				s: "This train has been delayed by a fault with the signalling system",
				l: 40,
			},
			want: []string{
				"This train has been delayed by a fault",
				"with the signalling system",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.args.s, tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}
