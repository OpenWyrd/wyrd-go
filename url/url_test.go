package url

import "testing"

const (
	validHandle = "AAAAAAAAAAAAAAAA"                              // 16 chars
	validKRead  = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"   // 43 chars
)

func TestParseFragment(t *testing.T) {
	got := Parse("https://sendwyrd.com/w/" + validHandle + "#" + validKRead)
	if got.Form != FormFragment {
		t.Fatalf("expected FormFragment, got %v (reason=%q)", got.Form, got.Reason)
	}
	if got.Handle != validHandle {
		t.Errorf("handle = %q, want %q", got.Handle, validHandle)
	}
	if got.KRead != validKRead {
		t.Errorf("k_read = %q, want %q", got.KRead, validKRead)
	}
}

func TestParsePublicLegacy(t *testing.T) {
	got := Parse("https://sendwyrd.com/w/" + validHandle + "/k/" + validKRead)
	if got.Form != FormPublic {
		t.Fatalf("expected FormPublic, got %v (reason=%q)", got.Form, got.Reason)
	}
	if got.Handle != validHandle || got.KRead != validKRead {
		t.Errorf("handle=%q kread=%q", got.Handle, got.KRead)
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		reason string
	}{
		{"no_url", "::not::a::url::", "not_a_url"},
		{"wrong_path", "https://sendwyrd.com/notwyrd/foo", "path_does_not_match_wyrd_url"},
		{"short_handle", "https://sendwyrd.com/w/AAAA#" + validKRead, "path_does_not_match_wyrd_url"},
		{"missing_fragment", "https://sendwyrd.com/w/" + validHandle, "missing_or_malformed_fragment"},
		{"short_fragment", "https://sendwyrd.com/w/" + validHandle + "#AAAA", "missing_or_malformed_fragment"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.input)
			if got.Form != FormInvalid {
				t.Fatalf("expected FormInvalid, got %v", got.Form)
			}
			if got.Reason != tt.reason {
				t.Errorf("reason = %q, want %q", got.Reason, tt.reason)
			}
		})
	}
}

func TestBuildFragment(t *testing.T) {
	got := BuildFragment("https://sendwyrd.com", validHandle, validKRead)
	want := "https://sendwyrd.com/w/" + validHandle + "#" + validKRead
	if got != want {
		t.Errorf("BuildFragment = %q, want %q", got, want)
	}
}

func TestRoundTrip(t *testing.T) {
	url := BuildFragment("https://relay.example", validHandle, validKRead)
	got := Parse(url)
	if got.Form != FormFragment || got.Handle != validHandle || got.KRead != validKRead {
		t.Errorf("round-trip failed: %+v", got)
	}
}

// TestParseConformanceVector exercises Parse against a vector formatted to
// match the canonical inputs from openwyrd/mop-conformance/vectors/v1/.
// When the conformance suite ships, this test imports the vector JSON
// directly — for now it inlines the values from
// mop-conformance/vectors/v1/00-canonical-envelope.json.
func TestParseConformanceVector(t *testing.T) {
	// Canonical vector inputs from openwyrd/mop-conformance v1 schema.
	// See vectors/v1/00-canonical-envelope.json.
	const handle = "AAAAAAAAAAAAAAAA"
	// k_read for this vector is [COMPUTE:] in the conformance file —
	// the reference impl populates the actual base64url value in CI.
	// Here we just verify the URL structure is consumed correctly given
	// any well-formed k_read.
	const kRead = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	cap := BuildFragment("https://sendwyrd.com", handle, kRead)
	parsed := Parse(cap)
	if parsed.Form != FormFragment {
		t.Fatalf("conformance round-trip failed: form=%v reason=%q", parsed.Form, parsed.Reason)
	}
	if parsed.Handle != handle || parsed.KRead != kRead {
		t.Errorf("conformance round-trip mismatch: got %+v", parsed)
	}
}
