// Package url parses and constructs OpenWyrd MOP capability URLs per
// MOP-001 §5.
//
// Canonical form (§5.1):
//
//	https://<relay>/w/{handle}#{K_read_b64u}
//
// Legacy path form (§5.2) is recognized at parse time only; conformant
// composers never emit it.
package url

import (
	"fmt"
	"net/url"
	"regexp"
)

// Form is the kind of MOP URL a parse produced.
type Form int

const (
	// FormInvalid indicates the input did not match any MOP URL form.
	FormInvalid Form = iota
	// FormFragment is the canonical form: /w/{handle}#{k_read}.
	FormFragment
	// FormPublic is the legacy path form: /w/{handle}/k/{k_read}.
	FormPublic
)

// Parsed holds the result of [Parse].
type Parsed struct {
	Form   Form
	Handle string // base64url, 16 chars
	KRead  string // base64url, 43 chars (only set when Form is Fragment or Public)
	Reason string // populated when Form == FormInvalid
}

const (
	handleChars = 16 // 12 raw bytes → 16 base64url chars (no padding)
	kReadChars  = 43 // 32 raw bytes → 43 base64url chars (no padding)
)

var (
	// /w/{handle} — fragment form path.
	handlePathRe = regexp.MustCompile(fmt.Sprintf(`^/w/([A-Za-z0-9_-]{%d})$`, handleChars))

	// /w/{handle}/k/{k_read} — legacy path form.
	pathFormRe = regexp.MustCompile(fmt.Sprintf(`^/w/([A-Za-z0-9_-]{%d})/k/([A-Za-z0-9_-]{%d})$`, handleChars, kReadChars))

	// strict base64url charset for fragment validation.
	b64URLRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

// Parse classifies an input URL string. The HTTP scheme and host are not
// validated here — relay-host validation belongs at the transport layer.
// MOP cares only about the path + fragment shape.
func Parse(input string) Parsed {
	u, err := url.Parse(input)
	if err != nil {
		return Parsed{Form: FormInvalid, Reason: "not_a_url"}
	}

	// Legacy path form first — if the URL has /k/ in path, it's
	// unambiguously path-form (no fragment expected).
	if m := pathFormRe.FindStringSubmatch(u.Path); m != nil {
		return Parsed{Form: FormPublic, Handle: m[1], KRead: m[2]}
	}

	if m := handlePathRe.FindStringSubmatch(u.Path); m != nil {
		// url.URL.Fragment is already the un-prefixed fragment (no leading #),
		// and it's URL-decoded. We expect raw base64url chars only.
		frag := u.Fragment
		if len(frag) == 0 || len(frag) != kReadChars {
			return Parsed{
				Form:   FormInvalid,
				Handle: m[1],
				Reason: "missing_or_malformed_fragment",
			}
		}
		if !b64URLRe.MatchString(frag) {
			return Parsed{
				Form:   FormInvalid,
				Handle: m[1],
				Reason: "fragment_not_base64url",
			}
		}
		return Parsed{Form: FormFragment, Handle: m[1], KRead: frag}
	}

	return Parsed{Form: FormInvalid, Reason: "path_does_not_match_wyrd_url"}
}

// BuildFragment composes a canonical fragment-form capability URL.
// Per MOP-001 §5.1 / ADR-021 in the SendWyrd reference, this is the only
// form a conformant composer emits.
func BuildFragment(origin, handle, kRead string) string {
	return fmt.Sprintf("%s/w/%s#%s", origin, handle, kRead)
}
