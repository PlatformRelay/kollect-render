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

## Threat model

kollect-render is a **local, credential-free, network-free** CLI. It does not open network
connections, hold cloud credentials, talk to Kubernetes, or publish to remote systems. Callers
supply every path; the process runs with the caller's OS privileges.

### Trust boundary

| Input | Trust assumption | Primary risks |
|-------|------------------|---------------|
| `--context` inventory documents | **Untrusted** unless the caller verified provenance | Malicious or oversized YAML/JSON; unexpected field shapes; data that becomes XSS when a downstream HTML renderer consumes markdown output |
| `--template` Go templates | **Untrusted** unless the caller controls the template store | Template injection (arbitrary filesystem reads via template functions, unexpected execution paths, denial of service via pathological templates) |
| `--upstream-deps` and similar side inputs | **Untrusted** | Same class as context documents |
| `--output` path | Caller-chosen write target | Path traversal / overwrite of unintended files (including `.meta.json` sidecars next to the body path) |

### What this tool does

- **Reads** inventory documents and templates from paths you provide.
- **Writes** render artifacts (and optional `.meta.json` sidecars) to paths you provide.
- **Does not** fetch URLs, embed secrets, or escape the process sandbox beyond normal file I/O.

### Mitigations callers should apply

- Treat context files and templates as untrusted input when they come from outside your control
  (shared repos, CI artifacts, user uploads). Validate or pin template sources.
- Prefer the built-in format paths when you do not need custom templates; review any custom
  `text/template` / `html/template` carefully before use.
- Point `--output` at a dedicated directory; avoid writing into shared system paths.
- For **confluence-storage**, inventory strings are contextually XML-escaped. For **markdown**,
  strings are passed through — sanitize separately if you later HTML-render the result.
- Prefer tagged release artifacts over development builds in production pipelines.

Report supply-chain concerns (compromised releases, unexpected binary behaviour) through the
private contact above.
