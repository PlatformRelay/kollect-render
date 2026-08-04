# Output formats

Formats register at compile time in `internal/format`. Today the registry exposes two names,
returned in stable order by `format.Names()`:

| Name | Encoder | Typical consumer |
| --- | --- | --- |
| `markdown` | Built-in Model → GitHub-flavoured markdown | PRs, Git sinks, human review |
| `confluence-storage` | Built-in Model → Confluence storage-format XHTML | Private Confluence publisher |

Both encoders consume the same format-agnostic `format.Model` (a sealed set of block and inline
types). There is no raw-HTML / unsafe bypass type on the model.

## Markdown

- Headings, banners, status legend, paragraphs, bullet lists, tables, cap notes, and footnotes.
- Optional `{#anchor}` suffix on headings when an anchor is set.
- Inventory strings (including hostile markup such as `<script>…</script>` in a component name)
  are **passed through** into markdown text. Treat markdown output as untrusted HTML if you later
  render it in a browser without a separate sanitizer.
- Optional `--template` runs Go `text/template` against the render context (markdown format only).

## Confluence storage

- Emits Confluence storage-format markup (`<p>`, `<ul>`, `<table>`, `ac:` macros for status and
  anchors).
- Every string is **contextually XML-escaped** (`&`, `<`, `>`, `"`) before emission. Hostile
  markup in inventory fields becomes `&lt;script&gt;…`, not executable tags.
- Output is validated as well-formed XHTML (Confluence prefixes declared on a wrapper so standard
  XML tokenization applies).
- `--template` is not supported with this format.

## Choosing a format

Use `markdown` when the artifact will live in Git or be reviewed as plain text. Use
`confluence-storage` when a private publisher will push the body into Confluence — escaping is
part of the publish-ready contract.
