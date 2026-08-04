#!/usr/bin/env bash
# Fail closed unless a release commit is on protected main with exact-SHA checks
# and a Keep-a-Changelog section for the release version. Read-only — never tags
# or publishes. Intended to run from a trusted checkout (default branch).
set -euo pipefail

SHA="${1:?usage: verify-eligibility.sh <full-commit-sha> <version>}"
VERSION="${2:?usage: verify-eligibility.sh <full-commit-sha> <version>}"
REPO="${GITHUB_REPOSITORY:?GITHUB_REPOSITORY required}"
DEFAULT_BRANCH="${RELEASE_DEFAULT_BRANCH:-main}"

if [[ ! "${SHA}" =~ ^[0-9a-f]{40}$ ]]; then
	echo "error: release SHA must be a full 40-character lowercase commit SHA" >&2
	exit 1
fi

if [[ ! "${VERSION}" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
	echo "error: version must be SemVer-like (e.g. 0.2.0 or 0.2.0-rc.1), without a leading v" >&2
	exit 1
fi

required_checks=(
	check
	changelog
)

main_sha="$(gh api "repos/${REPO}/commits/${DEFAULT_BRANCH}" --jq .sha)"
comparison="$(gh api "repos/${REPO}/compare/${SHA}...${main_sha}" --jq .status)"
if [[ "${comparison}" != "ahead" && "${comparison}" != "identical" ]]; then
	echo "error: release SHA ${SHA} is not reachable from protected main (${main_sha})" >&2
	exit 1
fi

# gh ≥2.x rejects combining --slurp with --jq; flatten pages in a separate jq pass.
checks_pages="$(gh api --paginate --slurp "repos/${REPO}/commits/${SHA}/check-runs?per_page=100")"
checks="$(jq '[.[].check_runs[] | select(.head_sha == "'"${SHA}"'") | {id, name, status, conclusion, head_sha}]' \
	<<<"${checks_pages}")"

failed=0
for name in "${required_checks[@]}"; do
	result="$(jq -r --arg name "${name}" '
		(map(select(.name == $name)) | sort_by(.id) | last) as $run |
		if $run == null then "missing"
		else ($run.status + "/" + ($run.conclusion // ""))
		end
	' <<<"${checks}")"
	if [[ "${result}" != "completed/success" ]]; then
		echo "error: required exact-SHA check ${name}: ${result}" >&2
		failed=1
	else
		echo "ok: ${name}"
	fi
done
if [[ "${failed}" -ne 0 ]]; then
	exit 1
fi

changelog="$(gh api \
	-H "Accept: application/vnd.github.raw" \
	"repos/${REPO}/contents/CHANGELOG.md?ref=${SHA}")"
if ! grep -Eq "^## \\[${VERSION}\\]( |$)" <<<"${changelog}"; then
	echo "error: CHANGELOG.md has no section for ${VERSION} at ${SHA}" >&2
	exit 1
fi

echo "Release eligibility passed for ${REPO}@${SHA} (CHANGELOG section [${VERSION}], main ${main_sha})."
