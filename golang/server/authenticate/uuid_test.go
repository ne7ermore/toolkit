package authenticate

import "testing"

func TestUUID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"asd", "aaa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := UUID(); got != tt.want {
				t.Errorf("UUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
