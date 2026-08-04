# Security

## Threat model (summary)

kollect-render is a local, credential-free document tool:

- **Reads** inventory documents and templates from paths you provide.
- **Writes** render artifacts (and optional `.meta.json` sidecars) to paths you provide.
- **Does not** open network connections, hold cloud credentials, or publish to remote systems.

Prefer tagged release artifacts over development builds in production pipelines.

## Output safety

- **confluence-storage** contextually XML-escapes all inventory strings before emission.
- **markdown** passes inventory strings through; if you later HTML-render markdown, sanitize
  separately.
- The format `Model` has no raw-HTML / unsafe bypass type — encoders own escaping.

## Supported versions

| Version | Supported |
| --- | --- |
| `main` | Yes |
| Tags | Latest release only |

## Reporting a vulnerability

**Do not** open public GitHub issues for security-sensitive reports.

Email **konrad.heimel@gmail.com** with:

- Description of the issue and impact
- Steps to reproduce (if possible)
- Affected versions or commits
- Suggested fix (optional)

You should receive an acknowledgment within a few business days. We will coordinate disclosure
and a fix release when appropriate. The same policy is mirrored in
[SECURITY.md](https://github.com/PlatformRelay/kollect-render/blob/main/SECURITY.md) at the
repository root.
