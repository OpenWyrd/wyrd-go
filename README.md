# wyrd-go

Go implementation of [OpenWyrd MOP](https://github.com/openwyrd/mop). Single static binary, SQLite + filesystem storage, runs on a $5 VPS or homelab via `docker-compose up`.

The point of this implementation is **interop**, not feature parity with the canonical TypeScript reference. If `wyrd-go` and the reference can round-trip a wyrd byte-for-byte through the conformance suite, the wire spec is real. If they diverge, the spec is wrong.

## Status

**Pre-alpha. Scaffolding.** No crypto, no server, no storage yet. See [`MIGRATION-PLAN.md`](MIGRATION-PLAN.md) for the ordered work.

What works today:
- Module layout
- Capability-URL parsing (`url/`) — the simplest MOP-001 §5 surface, no crypto needed
- CI: `go vet` + `go test ./...` on push

What doesn't:
- Everything else.

## Layout

```
wyrd-go/
├── cmd/wyrd-go/      # Server binary entry point (stub)
├── wire/             # Types matching MOP-001 §4 (envelope, AAD, request/response)
├── url/              # Capability URL construction + parsing (MOP-001 §5)
├── crypto/           # secp256k1 + ECIES + AES-256-GCM + HKDF-SHA256 (MOP-001 §6) — TBD
├── relay/            # HTTP server: publish, fetch, burn, replies (MOP-001 §8) — TBD
├── store/            # SQLite + filesystem persistence (MOP-001 §8.1) — TBD
└── MIGRATION-PLAN.md # Ordered work breakdown
```

## Why Go

`@noble/secp256k1` and `@noble/ciphers` are excellent — and reusing them across implementations would defeat the purpose. A Go port of the crypto forces every implicit assumption in the TypeScript reference to surface: SEC1 vs raw point encodings, JSON canonicalization (Go's `encoding/json` differs from V8's `JSON.stringify`), AES-GCM tag placement (`crypto/cipher` appends; some JS APIs return separately), base64url padding strictness, BIP-340 Schnorr nonce derivation. Every divergence the conformance suite catches is a wire-spec gap closed.

Go was picked over Rust for contribution barrier (Go is the lingua franca of self-hosted infra), over Python for deployment ergonomics ($5 VPS, no venv), over Bun/Deno because reusing JS crypto libraries defeats the whole interop test.

## Running (eventually)

```bash
go build -o wyrd-go ./cmd/wyrd-go
./wyrd-go --addr :8080 --db wyrd.db --storage ./blobs
```

## Conformance

Once crypto + relay are landed, this binary must pass the [`openwyrd/mop-conformance`](https://github.com/openwyrd/mop-conformance) suite at MOP-001@1.0.0 to claim compliance. CI runs the full vector set against every commit.

## License

This repo is **mixed-license** — see [`LICENSING.md`](LICENSING.md) for the full breakdown. Short version:

- **Library code** (`url/`, `wire/`, `crypto/`, `store/`) — Apache-2.0. Vendor it freely; the trust boundary is the server, not the libraries.
- **Daemon code** (`cmd/wyrd-go/`, `relay/`) — AGPL-3.0-or-later. The hosted-service loophole closes here; running a modified server against users obligates you to publish your modifications.

The split mirrors the JS side: [`@openwyrd/mop`](https://github.com/openwyrd/mop-js) is permissive, [`openwyrd/sendwyrd`](https://github.com/openwyrd/sendwyrd) is copyleft. Repository root [`LICENSE`](LICENSE) is AGPL-3.0-or-later as the default; each Apache-2.0 directory carries its own [`LICENSE`](url/LICENSE).
