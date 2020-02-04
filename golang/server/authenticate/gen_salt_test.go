package authenticate

import "testing"

func Test_getSalt(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{16}, "asdasd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSalt(tt.args.n); got != tt.want {
				t.Errorf("getSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}
