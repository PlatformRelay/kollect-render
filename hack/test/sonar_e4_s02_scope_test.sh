#!/usr/bin/env bash
# E4-S02: Sonar scope — CPD + go:S3776 must ignore tests only, not production.
# Arrange-invoke-assert clones in *_test.go are intentional CLI contract
# readability; Cognitive Complexity (S3776) in production stays open.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROPS="${ROOT}/sonar-project.properties"

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

pass() {
  echo "ok - $*"
}

[[ -f "${PROPS}" ]] || fail "sonar-project.properties not found at ${PROPS}"

# Production paths with known/possible S3776 must remain analyzed (not silenced).
for prod in \
  cmd/kollect-render/main.go \
  internal/format/env_inventory.go \
  internal/format/confluence.go
do
  [[ -f "${ROOT}/${prod}" ]] || fail "production file missing: ${prod}"
done
pass "production S3776 target files exist"

# CPD exclusions: tests only
cpd_line="$(grep -E '^sonar\.cpd\.exclusions=' "${PROPS}" || true)"
[[ -n "${cpd_line}" ]] || fail "sonar.cpd.exclusions= not found in ${PROPS}"
cpd_value="${cpd_line#sonar.cpd.exclusions=}"
[[ "${cpd_value}" == "**/*_test.go" ]] || fail "sonar.cpd.exclusions expected **/*_test.go, got: ${cpd_value}"
pass "sonar.cpd.exclusions=**/*_test.go"

# Multicriteria master key includes s3776tests
master_line="$(grep -E '^sonar\.issue\.ignore\.multicriteria=' "${PROPS}" || true)"
[[ -n "${master_line}" ]] || fail "sonar.issue.ignore.multicriteria= not found in ${PROPS}"
master_value="${master_line#sonar.issue.ignore.multicriteria=}"
IFS=',' read -r -a IDS <<<"${master_value}"
found=0
for id in "${IDS[@]}"; do
  id="${id//[[:space:]]/}"
  [[ "${id}" == "s3776tests" ]] && found=1
done
((found == 1)) || fail "multicriteria master key missing s3776tests (got: ${master_value})"
pass "multicriteria lists s3776tests"

rule_line="$(grep -E '^sonar\.issue\.ignore\.multicriteria\.s3776tests\.ruleKey=' "${PROPS}" || true)"
[[ -n "${rule_line}" ]] || fail "s3776tests.ruleKey missing"
[[ "${rule_line#*=}" == "go:S3776" ]] || fail "s3776tests.ruleKey expected go:S3776, got: ${rule_line#*=}"
pass "s3776tests.ruleKey=go:S3776"

res_line="$(grep -E '^sonar\.issue\.ignore\.multicriteria\.s3776tests\.resourceKey=' "${PROPS}" || true)"
[[ -n "${res_line}" ]] || fail "s3776tests.resourceKey missing"
[[ "${res_line#*=}" == "**/*_test.go" ]] || fail "s3776tests.resourceKey expected **/*_test.go, got: ${res_line#*=}"
pass "s3776tests.resourceKey=**/*_test.go"

# Must NOT exclude tests from analysis entirely (unlike kollect).
excl_line="$(grep -E '^sonar\.exclusions=' "${PROPS}" || true)"
if [[ -n "${excl_line}" ]]; then
  excl_value="${excl_line#sonar.exclusions=}"
  IFS=',' read -r -a EXCL <<<"${excl_value}"
  for pat in "${EXCL[@]}"; do
    pat="${pat//[[:space:]]/}"
    if [[ "${pat}" == "**/*_test.go" || "${pat}" == "*_test.go" ]]; then
      fail "sonar.exclusions must not drop tests from analysis (found: ${pat})"
    fi
  done
fi
pass "tests remain in sonar analysis set"

echo "All sonar_e4_s02 scope tests passed."
