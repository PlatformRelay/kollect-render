# Security policy

## Supported versions

| Version | Supported |
|---------|-----------|
| `main`  | Yes       |
| Tags    | Latest release only |

## Reporting a vulnerability

**Do not** open public GitHub issues for security-sensitive reports.

Email **konrad.heimel@gmail.com** with:

- Description of the issue and impact
- Steps to reproduce (if possible)
- Affected versions or commits
- Suggested fix (optional)

You should receive an acknowledgment within a few business days. We will coordinate disclosure
and a fix release when appropriate.

## Threat model (summary)

kollect-render is a local, credential-free document tool:

- **Reads** inventory documents and templates from paths you provide.
- **Writes** render artifacts to paths you provide.
- **Does not** open network connections, hold cloud credentials, or publish to remote systems.

Prefer tagged release artifacts over development builds in production pipelines. Report
supply-chain concerns through the private contact above.
