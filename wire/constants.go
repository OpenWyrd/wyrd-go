// Package wire defines the byte-level types of the OpenWyrd MOP wire format.
// See https://github.com/openwyrd/mop/blob/main/spec/MOP-001.md §4.
package wire

// Constants frozen by MOP-001 §4.1.
const (
	// HandleBytes is the length of the per-wyrd identifier in bytes.
	// Encodes to 16 ASCII characters in base64url (no padding).
	HandleBytes = 12

	// EnvelopeVersion is the leading byte of a v1 envelope. See §4.2.
	EnvelopeVersion byte = 0x01

	// AADVersion is the leading byte of the AAD. See §4.3.
	AADVersion byte = 0x01

	// IVBytes is the AES-256-GCM IV length.
	IVBytes = 12

	// TagBytes is the AES-256-GCM authentication tag length.
	TagBytes = 16

	// PermanentExpiresAtMs is the sentinel expiry timestamp for ttl=0
	// (permanent) wyrds. The AAD binds to this constant rather than to a
	// live "now + 0" value. See MOP-001 §7.2 — the permanent-row footgun.
	//
	// Year 9999 (= 253370764800 seconds since epoch, * 1000 for ms).
	PermanentExpiresAtMs int64 = 253_370_764_800_000

	// MaxBodyCodepoints is the body codepoint cap (§4.1).
	MaxBodyCodepoints = 300

	// MaxTTLSeconds is the maximum non-permanent TTL (§7.1).
	MaxTTLSeconds = 31_536_000 // one year

	// TombstoneRetentionDays is the retention window for 410 tombstones
	// (§8.1).
	TombstoneRetentionDays = 30
)

// HKDFInfoKRead is the v1-frozen HKDF info prefix for K_read derivation.
// MOP-001 §6.2 freezes this exact byte string for v1 conformance — do
// NOT "correct" to "mop:k_read" without bumping the spec version.
const HKDFInfoKRead = "sendwyrd:k_read"

// SignDomain is the v1-frozen prefix for signed-message digests.
// MOP-001 §6.4. Values: publish, delete, fetch_replies, presence_check.
var SignDomain = struct {
	Publish        string
	Delete         string
	FetchReplies   string
	PresenceCheck  string
}{
	Publish:       "mop:v1:publish",
	Delete:        "mop:v1:delete",
	FetchReplies:  "mop:v1:fetch_replies",
	PresenceCheck: "mop:v1:presence_check",
}

// ReplyHKDF is the v1-frozen ECIES HKDF info prefix for one-shot replies.
// MOP-001 §6.5.
var ReplyHKDF = struct {
	AESKey string
	IV     string
}{
	AESKey: "mop:v1:reply:aes_key:",
	IV:     "mop:v1:reply:iv:",
}
