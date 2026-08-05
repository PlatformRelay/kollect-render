#!/usr/bin/env bash
# E4-S02: Sonar scope — CPD must not count *_test.go; go:S3776 ignore tests only.
#
# SonarCloud indexes files matching both sonar.sources=. and
# sonar.test.inclusions as sources for CPD. sonar.cpd.exclusions alone does
# not clear new_duplicated_lines_density when tests are dual-indexed (QG
# evidence: all duplicated files were *_test.go despite cpd.exclusions).
# Standard split: exclude *_test.go from sources, keep them as tests so
# issue rules still apply; S3776 stays scoped via multicriteria.
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

# Sources must exclude tests so CPD does not dual-index them as main sources.
excl_line="$(grep -E '^sonar\.exclusions=' "${PROPS}" || true)"
[[ -n "${excl_line}" ]] || fail "sonar.exclusions= not found in ${PROPS}"
excl_value="${excl_line#sonar.exclusions=}"
IFS=',' read -r -a EXCL <<<"${excl_value}"
found_test_excl=0
for pat in "${EXCL[@]}"; do
  pat="${pat//[[:space:]]/}"
  if [[ "${pat}" == "**/*_test.go" ]]; then
    found_test_excl=1
  fi
done
((found_test_excl == 1)) || fail "sonar.exclusions must include **/*_test.go (got: ${excl_value})"
pass "sonar.exclusions includes **/*_test.go"

# Tests remain analyzed as tests (issues / coverage linkage), not dropped entirely.
test_inc_line="$(grep -E '^sonar\.test\.inclusions=' "${PROPS}" || true)"
[[ -n "${test_inc_line}" ]] || fail "sonar.test.inclusions= not found in ${PROPS}"
[[ "${test_inc_line#sonar.test.inclusions=}" == "**/*_test.go" ]] || \
  fail "sonar.test.inclusions expected **/*_test.go, got: ${test_inc_line#*=}"
pass "sonar.test.inclusions=**/*_test.go"

# Belt-and-suspenders: keep cpd.exclusions aligned (effective fix is source exclusion).
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

# Production must not be blanket-excluded from analysis.
for prod_pat in \
  '**/cmd/**' \
  '**/internal/**' \
  'cmd/kollect-render/main.go' \
  'internal/format/**'
do
  for pat in "${EXCL[@]}"; do
    pat="${pat//[[:space:]]/}"
    if [[ "${pat}" == "${prod_pat}" ]]; then
      fail "sonar.exclusions must not silence production path ${prod_pat}"
    fi
  done
done
pass "production paths not blanket-excluded"

echo "All sonar_e4_s02 scope tests passed."
