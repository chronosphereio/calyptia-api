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
			name:        "test empty project token",
			token:       []byte(""),
			projectID:   "",
			expectError: true,
		},
		{
			name:        "test malformed project token",
			token:       []byte("a.b"),
			projectID:   "",
			expectError: true,
		},

		{
			name:        "test invalid project token",
			token:       []byte("asdf"),
			projectID:   "",
			expectError: true,
		},
		{
			name:        "invalid project token",
			token:       []byte("ayJUb2tlbklEIjoiOWI5ODIwNzUtYzY4MC00MzdiLWE4YjMtYjU5NjNkMzE4OTUyIiwiUHJvamVjdElEIjoiMDkwZDFhYTEtZGU5Ni00NDZjLTk1NDQtMGUwMGNiNmRkMzkzIn1.JPD_g6oDQdmO_sPlshUJdNefpHT7AMDUjSRjg0x0E61U8-Frh2_ZOCTP93O5UBC9"),
			projectID:   "",
			expectError: true,
		},
		{
			name:        "test valid project token",
			token:       []byte("eyJUb2tlbklEIjoiOWI5ODIwNzUtYzY4MC00MzdiLWE4YjMtYjU5NjNkMzE4OTUyIiwiUHJvamVjdElEIjoiMDkwZDFhYTEtZGU5Ni00NDZjLTk1NDQtMGUwMGNiNmRkMzkzIn0.JPD_g6oDQdmO_sPlshUJdNefpHT7AMDUjSRjg0x0E61U8-Frh2_ZOCTP93O5UBC9"),
			projectID:   "090d1aa1-de96-446c-9544-0e00cb6dd393",
			expectError: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Decode(tc.token)
			if err != nil {
				if !tc.expectError {
					t.Errorf("error: %v != nil", err)
				}

				return
			}
			if want, got := tc.projectID, d; want != got {
				t.Errorf("want: %v != got: %v", want, got)
				return
			}
		})
	}
}
