# eKasa Go Client Project Instructions

## Architecture
- **Single Package**: All core logic, models, and constants are in the root `ekasa` package for simplicity and ease of use.
- **Functional Options**: Configuration follows the functional options pattern in `options.go`.
- **NWS4 Signing**: Implemented in `signer.go` to match the custom AWS-like signing used by Nine Digit APIs. This includes canonicalizing requests and computing HMAC-SHA256 signatures.

## Conventions
- **Pointers for Nullability**: Struct fields that can be null in JSON use pointers (e.g., `*time.Time`, `*string`, `*int`) to correctly handle missing or null values during JSON marshalling/unmarshalling.
- **Enums as Constants**: Enum-like types are implemented as `const` blocks with appropriate prefixes (e.g., `RegistrationState...`, `ReceiptType...`).
- **Context Support**: All API methods accept a `context.Context` as the first argument to support timeouts, cancellations, and tracing.

## Documentation
- Detailed usage, configuration examples, and environment details are located in the `docs/` directory.

## Maintenance
- Use `go fmt` and `go vet` to maintain code quality.
- The `signer.go` implementation is sensitive to formatting (canonicalization). When modifying, ensure it remains compatible with the PHP reference implementation.
