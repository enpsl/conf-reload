package errors

import "testing"

func TestErrorIs(t *testing.T) {
	var ErrCustom = New("custom error")

	tests := []struct {
		desc   string
		err    error
		target error
		want   bool
	}{
		{
			desc:   "errType is equal",
			err:    ErrFormat(ErrInvalidFilePath, ErrCustom),
			target: ErrInvalidFilePath,
			want:   true,
		},
	}

	for _, tc := range tests {
		if got := Is(tc.err, tc.target); got != tc.want {
			t.Errorf("%s: got=%t, want=%t", tc.desc, got, tc.want)
		}
	}
}

func TestErrorAs(t *testing.T) {
	var ErrCustom = New("custom error")
	tests := []struct {
		desc   string
		err    error
		target interface{}
		want   bool
	}{
		{
			desc:   "error type is equal",
			err:    ErrFormat(ErrInvalidFilePath, ErrCustom),
			target: ErrCustom,
			want:   true,
		},
	}

	for _, tc := range tests {
		if got := As(tc.err, &tc.target); got != tc.want {
			t.Errorf("%s: got=%t, want=%t", tc.desc, got, tc.want)
		}
	}
}
