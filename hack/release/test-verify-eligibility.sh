#!/usr/bin/env bash
# Deterministic fixture tests for verify-eligibility.sh (no live GitHub API).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SUBJECT="${ROOT}/hack/release/verify-eligibility.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "${TMP}"' EXIT

SHA="1111111111111111111111111111111111111111"
MAIN="2222222222222222222222222222222222222222"
VERSION="0.2.0"

REQUIRED_NAMES='check changelog'

cat >"${TMP}/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
args="$*"
# Mirror gh ≥2.x: --slurp cannot be combined with --jq/--template.
if [[ "${args}" == *--slurp* && ( "${args}" == *--jq* || "${args}" == *--template* ) ]]; then
	echo 'the `--slurp` option is not supported with `--jq` or `--template`' >&2
	exit 1
fi
case "${args}" in
  *"/commits/main"*)
    printf '%s\n' "${MOCK_MAIN}"
    ;;
  *"/compare/"*)
    printf '%s\n' "${MOCK_COMPARE}"
    ;;
  *"/check-runs"*)
    cat "${MOCK_CHECKS}"
    ;;
  *"/contents/CHANGELOG.md"*)
    cat "${MOCK_CHANGELOG}"
    ;;
  *)
    echo "unexpected gh invocation: ${args}" >&2
    exit 70
    ;;
esac
MOCK
chmod +x "${TMP}/gh"

write_green_checks() {
	# Shape matches `gh api --paginate --slurp` (array of pages).
	jq -n \
		--arg names "${REQUIRED_NAMES}" \
		--arg sha "${SHA}" \
		'[{check_runs: ($names | split(" ") | to_entries | map({
			id: (.key + 1),
			name: .value,
			status: "completed",
			conclusion: "success",
			head_sha: $sha
		}))}]' >"${TMP}/checks.json"
}

write_changelog() {
	local body="$1"
	# Raw contents Accept header returns plain text (script uses -H Accept: raw).
	printf '%s\n' "${body}" >"${TMP}/changelog.txt"
}

run_case() {
	PATH="${TMP}:${PATH}" \
		GITHUB_REPOSITORY=PlatformRelay/kollect-render \
		MOCK_MAIN="${MAIN}" \
		MOCK_COMPARE="${MOCK_COMPARE}" \
		MOCK_CHECKS="${TMP}/checks.json" \
		MOCK_CHANGELOG="${TMP}/changelog.txt" \
		bash "${SUBJECT}" "${SHA}" "${VERSION}" >"${TMP}/out" 2>"${TMP}/err"
}

fail_if_passes() {
	local label="$1"
	if run_case; then
		echo "${label} unexpectedly passed" >&2
		exit 1
	fi
}

write_green_checks
write_changelog "## [Unreleased]

## [${VERSION}] - 2026-08-04

### Features

- Something releasable
"

if [[ ! -x "${SUBJECT}" && ! -f "${SUBJECT}" ]]; then
	echo "missing subject ${SUBJECT}" >&2
	exit 1
fi

if PATH="${TMP}:${PATH}" GITHUB_REPOSITORY=PlatformRelay/kollect-render \
	bash "${SUBJECT}" bad-sha "${VERSION}" >"${TMP}/out" 2>"${TMP}/err"; then
	echo "invalid SHA unexpectedly passed" >&2
	exit 1
fi
grep -q 'full 40-character lowercase commit SHA' "${TMP}/err"

MOCK_COMPARE=diverged
fail_if_passes "non-main SHA"
grep -q 'not reachable from protected main' "${TMP}/err"

MOCK_COMPARE=ahead
jq '.[0].check_runs |= map(select(.name != "check"))' "${TMP}/checks.json" >"${TMP}/missing.json"
mv "${TMP}/missing.json" "${TMP}/checks.json"
fail_if_passes "missing check"
grep -q 'required exact-SHA check check: missing' "${TMP}/err"

write_green_checks
jq '(.[0].check_runs[] | select(.name == "check")).conclusion = "failure"' "${TMP}/checks.json" >"${TMP}/red.json"
mv "${TMP}/red.json" "${TMP}/checks.json"
fail_if_passes "red check"
grep -q 'required exact-SHA check check: completed/failure' "${TMP}/err"

write_green_checks
write_changelog "## [Unreleased]

### Features

- No release section yet
"
fail_if_passes "missing changelog"
grep -q "CHANGELOG.md has no section for ${VERSION}" "${TMP}/err"

write_green_checks
write_changelog "## [Unreleased]

## [${VERSION}] - 2026-08-04

### Features

- Something releasable
"
run_case
grep -q "Release eligibility passed for PlatformRelay/kollect-render@${SHA}" "${TMP}/out"
grep -q "CHANGELOG section \\[${VERSION}\\]" "${TMP}/out"

echo "verify-eligibility tests: ok"
