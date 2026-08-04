#!/usr/bin/env bash
# E5-S07: every Go package has a useful package comment; schema explains InventoryV0.
# Run: bash hack/test/pkgdocs_e5_s07_test.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

fail=0
ok() { echo "ok: $*"; }
bad() { echo "FAIL: $*" >&2; fail=1; }

# Extract the package-comment prose from `go doc` (before the first symbol listing).
package_prose() {
  local doc="$1"
  # Drop the "package name // import …" header line(s), keep text until a symbol line.
  awk '
    BEGIN { skip=1 }
    skip && /^package / { next }
    skip && /^[[:space:]]*$/ { next }
    skip { skip=0 }
    /^(const|var|func|type) / { exit }
    { print }
  ' <<<"$doc"
}

# Every listed package must print a "Package <name> …" synopsis via go doc.
while IFS= read -r pkg; do
  name="$(go list -f '{{.Name}}' "$pkg")"
  doc="$(go doc "$pkg" 2>&1 || true)"
  if [[ -z "$doc" ]]; then
    bad "go doc ${pkg} produced no output"
    continue
  fi
  prose="$(package_prose "$doc")"
  if ! grep -Eq "^Package ${name} " <<<"$prose"; then
    bad "go doc ${pkg} missing package comment (expected line starting with 'Package ${name} ')"
  else
    ok "package comment for ${pkg} (${name})"
  fi
done < <(go list ./...)

# schema package doc must explain InventoryV0 and the pinned draft-v0 schema file.
schema_doc="$(go doc github.com/platformrelay/kollect-render/schema 2>&1 || true)"
schema_prose="$(package_prose "$schema_doc")"
if grep -Eq 'InventoryV0' <<<"$schema_prose"; then
  ok "schema package doc mentions InventoryV0"
else
  bad "schema package doc must explain InventoryV0 (in package comment, not only the symbol list)"
fi
if grep -Eq 'inventory-v0\.schema\.json' <<<"$schema_prose"; then
  ok "schema package doc pins inventory-v0.schema.json"
else
  bad "schema package doc must name inventory-v0.schema.json"
fi
if grep -Eiq 'draft[[:space:]]*v0' <<<"$schema_prose"; then
  ok "schema package doc mentions draft v0"
else
  bad "schema package doc must mention draft v0"
fi

# Exported InventoryV0 must keep its own doc comment (revive exported).
inv_doc="$(go doc github.com/platformrelay/kollect-render/schema.InventoryV0 2>&1 || true)"
if grep -Eq 'InventoryV0' <<<"$inv_doc" && grep -Eiq 'schema|envelope' <<<"$inv_doc"; then
  ok "InventoryV0 exported doc present"
else
  bad "schema.InventoryV0 must have an exported doc comment"
fi

# Deliberate: no pkg.go.dev / Go Reference badge (everything meaningful is internal/).
readme="${ROOT}/README.md"
if grep -Eiq 'pkg\.go\.dev|Go[[:space:]]+Reference' "$readme"; then
  bad "README must not advertise a pkg.go.dev / Go Reference badge"
else
  ok "README has no pkg.go.dev badge"
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "pkgdocs_e5_s07_test: RED" >&2
  exit 1
fi
echo "pkgdocs_e5_s07_test: GREEN"
