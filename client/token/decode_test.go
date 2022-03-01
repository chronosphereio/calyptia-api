package token

import (
	"testing"
)

func TestDecode(t *testing.T) {
	tt := []struct {
		name        string
		token       []byte
		projectID   string
		expectError bool
	}{
		{
			"test empty project token",
			[]byte(""),
			"",
			true,
		},
		{
			"test malformed project token",
			[]byte("a.b"),
			"",
			true,
		},

		{
			"test invalid project token",
			[]byte("asdf"),
			"",
			true,
		},
		{
			"invalid project token",
			[]byte("ayJUb2tlbklEIjoiOWI5ODIwNzUtYzY4MC00MzdiLWE4YjMtYjU5NjNkMzE4OTUyIiwiUHJvamVjdElEIjoiMDkwZDFhYTEtZGU5Ni00NDZjLTk1NDQtMGUwMGNiNmRkMzkzIn1.JPD_g6oDQdmO_sPlshUJdNefpHT7AMDUjSRjg0x0E61U8-Frh2_ZOCTP93O5UBC9"),
			"",
			true,
		},
		{
			"test valid project token",
			[]byte("eyJUb2tlbklEIjoiOWI5ODIwNzUtYzY4MC00MzdiLWE4YjMtYjU5NjNkMzE4OTUyIiwiUHJvamVjdElEIjoiMDkwZDFhYTEtZGU5Ni00NDZjLTk1NDQtMGUwMGNiNmRkMzkzIn0.JPD_g6oDQdmO_sPlshUJdNefpHT7AMDUjSRjg0x0E61U8-Frh2_ZOCTP93O5UBC9"),
			"090d1aa1-de96-446c-9544-0e00cb6dd393",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Decode(tc.token)
			if err != nil {
				if !tc.expectError {
					t.Errorf("error: %v != nil", err)
				}
			}
			if want, got := tc.projectID, d; want != got {
				t.Errorf("want: %v != got: %v", want, got)
				return
			}
		})
	}
}
