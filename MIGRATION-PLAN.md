# wyrd-go — migration plan

Ordered work breakdown to take this from scaffolding to MOP-compliant. Each step is small enough to land in one session and produces a testable artifact.

Reference at every step: [MOP-001](https://github.com/openwyrd/mop/blob/main/spec/MOP-001.md) (the spec) and [`openwyrd/mop-conformance`](https://github.com/openwyrd/mop-conformance) (the vectors that gate compliance).

## Step 1 — wire types (MOP-001 §4)

Port the Go structs:

- `wire.Envelope` — version byte + IV + ciphertext + tag, byte layout per §4.2
- `wire.AAD` — version + handle + expires_at_be + replies_enabled, layout per §4.3
- `wire.PublishRequest`, `wire.FetchResponse`, `wire.Tombstone`, `wire.DeleteRequest`, `wire.ReplyEnvelope`
- `wire.Constants` — port the `PERMANENT_EXPIRES_AT_MS` and `TOMBSTONE_RETENTION_DAYS` literals from MOP-001 §4.1

No crypto yet. JSON encoding via `encoding/json` with explicit field tags for canonical names. Tests against the conformance vector inputs (the JSON layout, not the crypto outputs).

**Done when:** `go test ./wire/` passes parsing the canonical envelope vector's JSON inputs.

## Step 2 — base64url + canonicalization helpers

- `wire.EncodeB64URL` / `wire.DecodeB64URL` — RFC 4648 §5, no padding, strict
- `wire.CanonicalJSON` — RFC 8785 / JCS-compliant. Go's `encoding/json` defaults to non-canonical (HTML escaping, lexicographic ordering not guaranteed for maps). Use a JCS lib or hand-roll the small subset needed.

**Done when:** vectors in `mop-conformance/vectors/v1/*` for `aad/canonicalization` and base64url edge cases all pass.

## Step 3 — secp256k1 crypto (MOP-001 §6.1, §6.4)

- Adopt `github.com/decred/dcrd/dcrec/secp256k1/v4`
- `crypto.SchnorrSign(msg, priv) → sig` (BIP-340)
- `crypto.SchnorrVerify(msg, pub, sig)`
- `crypto.PublishMessage(handle, env, ttl, replies, ts) → 32-byte digest` per §6.4
- `crypto.DeleteMessage(handle, ts)`, `FetchRepliesMessage(handle, ts)`, `PresenceCheckMessage(pub, ts)`
- Domain prefix strings frozen as per MOP-001 §6.4 (`"mop:v1:publish"` etc — see the §6.5 freeze note for the mixed prefix convention)

**Done when:** `signature/*` vectors in `mop-conformance` all round-trip — the Go signer produces signatures the reference TS verifier accepts, and vice versa.

## Step 4 — AES-256-GCM body encryption (MOP-001 §6.3)

- `crypto.EncryptEnvelope(k_read, iv, plaintext, aad) → envelope` per §6.3
- `crypto.DecryptEnvelope(k_read, envelope, aad) → plaintext`
- `0x01` version byte prefix, IV concatenation, tag handling matching `@noble/ciphers` output

**Watch for:** Go's `crypto/cipher` appends the GCM tag to the ciphertext; `@noble/ciphers` returns it as a separate field by default. The wire format puts them adjacent, so the layouts converge — but the in-memory APIs differ. Document the convergence point with a test.

**Done when:** `envelope/round-trip` vectors pass.

## Step 5 — HD derivation (MOP-001 §6.2)

- `crypto.DeriveKOriginKeypair(seed, n) → (priv, pub)` via BIP-32 `m/300'/n'`
- `crypto.DeriveKRead(seed, n) → 32 bytes` via HKDF-SHA256 with info `"sendwyrd:k_read" || n_be4` (frozen string per spec §6.2)

Use `github.com/tyler-smith/go-bip32` and `github.com/tyler-smith/go-bip39`, or audit-friendly equivalents.

**Watch for:** the HKDF info string is `"sendwyrd:..."` not `"mop:..."` per the v1 freeze in MOP-001 §6.2. Easy to "correct" and break interop.

**Done when:** vectors derived from a fixed BIP-39 seed produce identical K_origin pubkeys and K_read bytes across both implementations.

## Step 6 — capability-URL ECIES reply (MOP-001 §6.5)

- `crypto.EncryptReply(k_origin_pub, plaintext)` — ephemeral keypair, ECDH, HKDF info strings `"mop:v1:reply:aes_key:"` and `"mop:v1:reply:iv:"` per spec, AES-256-GCM
- `crypto.DecryptReply(k_origin_priv, blob)`

**Done when:** `signature/authorship-attestation` and reply-ECIES vectors round-trip.

## Step 7 — store layer

- SQLite schema mirroring SendWyrd's D1 schema (handles, envelopes, expires_at, replies, etc.)
- Filesystem blob storage for envelopes (or stash in SQLite as BLOB if simpler)
- Tombstone retention window per MOP-001 §8.1

**Done when:** integration test publishes and fetches from disk.

## Step 8 — relay HTTP server (MOP-001 §8)

- `chi` or `net/http` for the surface
- Routes per MOP-001 §8.2-8.5: POST /api/v1/wyrds, GET /api/v1/wyrds/:handle, DELETE /api/v1/wyrds/:handle, POST /api/v1/wyrds/:handle/replies, GET (presence-check)
- CORS open by spec — any client can read from any host
- Rate limiting and auth follow MOP-001 §8 normative requirements

**Done when:** the canonical reference (sendwyrd.com) and a local wyrd-go binary can publish to either and fetch from the other in a Playwright-driven cross-host test.

## Step 9 — conformance CLI

- Implement the runner contract from `openwyrd/mop-conformance/runner/README.md`:
  - `wyrd-go conformance encode --vector <file>`
  - `wyrd-go conformance decode --vector <file>`
  - JSON output, exit codes per the contract

**Done when:** `wyrd-go conformance run vectors/v1/` passes 100% of the v1 vector set.

## Step 10 — Docker + deploy story

- `Dockerfile` (scratch-based or `gcr.io/distroless/static`)
- `docker-compose.yml` reference deployment
- Deploy guide: `https://docs.openwyrd.org/deploy-wyrd-go`

**Done when:** a fresh user can `docker compose up` and have a MOP-compliant relay running on their box.

---

## Non-goals (not in scope for `wyrd-go`)

- Web UI / native client (clients use `@openwyrd/mop` directly; relay is HTTP)
- HD wallet UX, mnemonic backup flows (client concern)
- Cloudflare-specific features (Durable Objects, R2, Workers KV) — explicitly avoided
- Performance optimization beyond "sane defaults" — interop > throughput at v1
