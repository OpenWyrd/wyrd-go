# Licensing

`wyrd-go` is a mixed-license repository. The split mirrors the project-wide rule: **library code is permissive, server/daemon code is copyleft.** The trust boundary is the server, not the libraries.

## Per-directory licensing

| Directory | License | Why |
|---|---|---|
| `url/` | Apache-2.0 | Capability-URL parsing — needed by clients, libraries, tooling, embedded contexts. Permissive so the Go ecosystem can adopt it without copyleft viral concern. |
| `wire/` | Apache-2.0 | Wire-format types and constants. Same logic as `url/`. |
| `crypto/` | Apache-2.0 | secp256k1 / ECIES / AES-256-GCM / HKDF primitives. Permissive so security-sensitive consumers — wallets, signers, audit tooling — can vendor it. |
| `store/` | Apache-2.0 | Storage interface contract. Implementations may be AGPL when deployed inside a server. |
| `cmd/wyrd-go/` | AGPL-3.0-or-later | The server binary entry point. The hosted-service loophole closes here. |
| `relay/` | AGPL-3.0-or-later | HTTP relay handlers — server-side. AGPL §13 attaches when this code is offered as a network service. |

The repository root [`LICENSE`](LICENSE) is AGPL-3.0-or-later — the default license for any directory not explicitly carved out above. Each Apache-2.0 directory carries its own [`LICENSE`](url/LICENSE) file containing the full Apache-2.0 text.

## SPDX headers

Every source file declares its license via an `SPDX-License-Identifier` comment near the top of the file:

```go
// SPDX-License-Identifier: Apache-2.0
```

or

```go
// SPDX-License-Identifier: AGPL-3.0-or-later
```

This is the canonical signal for license scanners, SBOM tooling (CycloneDX, SPDX), and CRA-style regulatory inventories. New files in this repo must carry one. CI will eventually enforce.

## Practical implications

- **Building a Go library that imports `github.com/openwyrd/wyrd-go/url` or `/wire` or `/crypto` or `/store`:** Apache-2.0 applies. Vendor freely, including in proprietary code, subject to Apache attribution requirements.
- **Building a fork of the server:** AGPL-3.0-or-later applies. If you run a modified server against users you do not control, you must publish your modifications.
- **Mixing:** Importing an Apache-2.0 library directory from inside an AGPL-3.0 file is fine (Apache → AGPL combination is permitted). The reverse — an Apache-2.0 file importing AGPL-internal symbols — is not. The directory boundaries enforce the direction.

## Why not all-AGPL or all-Apache

All-AGPL would chill Go-ecosystem adoption of the libraries. The Go community is allergic to GPL-family viral copyleft, often as a matter of policy rather than law, and `wyrd-go` is meant to be vendored as much as it is meant to be deployed.

All-Apache would forfeit the hosted-service capture protection that AGPL §13 exists to provide. The deployable binary needs the copyleft hook.

The split keeps the libraries embeddable and the daemon copyleft. This is the same shape the JS side uses — `@openwyrd/mop` (Apache) library + `sendwyrd` (AGPL) server.

## Plugin boundary

Plugins, sidecars, and external integrations that interact with the AGPL daemon only via documented APIs (HTTP, capability URLs, IPC) are separate works for AGPL §5 purposes. See [`openwyrd/mop` `governance/CHARTER.md` §Plugin boundary](https://github.com/openwyrd/mop/blob/main/governance/CHARTER.md).
