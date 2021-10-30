package justtest

import "testing"

func TestAddUpTwoNumbers(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "correct result",
			args: args{
				a: 1,
				b: 2,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddUpTwoNumbers(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("AddUpTwoNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
