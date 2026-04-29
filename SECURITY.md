# Security policy

## Reporting a vulnerability

Email `security@openwyrd.org` with details. We are pre-launch and do not yet publish a maintainer PGP key; treat the channel accordingly.

A maintainer PGP key will be published here once it exists. Cross-publishable at <https://keys.openpgp.org> and <https://openwyrd.org/.well-known/maintainer.asc>.

## Disclosure window

90 days from initial report. We will:

- Acknowledge receipt within 72 hours.
- Confirm or refute the issue within 14 days.
- Coordinate a fix and a public disclosure date.
- Credit the reporter in the disclosure unless they request otherwise.

If we cannot ship a fix within 90 days, we will say so publicly and explain why. Indefinite embargoes are not on the table.

## In scope

- Cryptographic correctness of this implementation against MOP-001.
- Wire-format compliance failures discoverable by the conformance suite.
- Memory safety / authentication bypass / data leak in the relay HTTP surface (once it lands).

## Not in scope

- Issues in upstream Go crypto libs (`dcrd/secp256k1`, `crypto/cipher`) — file with the upstream.
- Operational issues with hosted deployments running `wyrd-go` — those belong to the host operator.
- Self-inflicted attacks (a recipient leaks their capability URL).

## No bug bounty

Not yet. We will not waste your time pretending to run a bounty program before we have an audited codebase.
