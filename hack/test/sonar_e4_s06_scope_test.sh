#!/usr/bin/env bash
# E4-S06: Sonar scope — go:S1186 and godre:S8196 ignored only on model.go.
# Marker methods seal Block/Inline; -er renames would mislead. Other packages
# must still surface empty-method / interface-naming findings.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROPS="${ROOT}/sonar-project.properties"
MODEL_RES='**/internal/format/model.go'

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

pass() {
  echo "ok - $*"
}

[[ -f "${PROPS}" ]] || fail "sonar-project.properties not found at ${PROPS}"
[[ -f "${ROOT}/internal/format/model.go" ]] || fail "model.go missing"

# Multicriteria master key must list both E4-S06 ids (and keep s3776tests).
master_line="$(grep -E '^sonar\.issue\.ignore\.multicriteria=' "${PROPS}" || true)"
[[ -n "${master_line}" ]] || fail "sonar.issue.ignore.multicriteria= not found"
master_value="${master_line#sonar.issue.ignore.multicriteria=}"
IFS=',' read -r -a IDS <<<"${master_value}"
found_s1186=0
found_s8196=0
found_s3776=0
for id in "${IDS[@]}"; do
  id="${id//[[:space:]]/}"
  [[ "${id}" == "s1186model" ]] && found_s1186=1
  [[ "${id}" == "s8196model" ]] && found_s8196=1
  [[ "${id}" == "s3776tests" ]] && found_s3776=1
done
((found_s3776 == 1)) || fail "multicriteria master missing s3776tests (got: ${master_value})"
((found_s1186 == 1)) || fail "multicriteria master missing s1186model (got: ${master_value})"
((found_s8196 == 1)) || fail "multicriteria master missing s8196model (got: ${master_value})"
pass "multicriteria lists s3776tests,s1186model,s8196model"

rule_s1186="$(grep -E '^sonar\.issue\.ignore\.multicriteria\.s1186model\.ruleKey=' "${PROPS}" || true)"
[[ -n "${rule_s1186}" ]] || fail "s1186model.ruleKey missing"
[[ "${rule_s1186#*=}" == "go:S1186" ]] || fail "s1186model.ruleKey expected go:S1186, got: ${rule_s1186#*=}"
pass "s1186model.ruleKey=go:S1186"

res_s1186="$(grep -E '^sonar\.issue\.ignore\.multicriteria\.s1186model\.resourceKey=' "${PROPS}" || true)"
[[ -n "${res_s1186}" ]] || fail "s1186model.resourceKey missing"
[[ "${res_s1186#*=}" == "${MODEL_RES}" ]] || fail "s1186model.resourceKey expected ${MODEL_RES}, got: ${res_s1186#*=}"
pass "s1186model.resourceKey=${MODEL_RES}"

rule_s8196="$(grep -E '^sonar\.issue\.ignore\.multicriteria\.s8196model\.ruleKey=' "${PROPS}" || true)"
[[ -n "${rule_s8196}" ]] || fail "s8196model.ruleKey missing"
[[ "${rule_s8196#*=}" == "godre:S8196" ]] || fail "s8196model.ruleKey expected godre:S8196, got: ${rule_s8196#*=}"
pass "s8196model.ruleKey=godre:S8196"

res_s8196="$(grep -E '^sonar\.issue\.ignore\.multicriteria\.s8196model\.resourceKey=' "${PROPS}" || true)"
[[ -n "${res_s8196}" ]] || fail "s8196model.resourceKey missing"
[[ "${res_s8196#*=}" == "${MODEL_RES}" ]] || fail "s8196model.resourceKey expected ${MODEL_RES}, got: ${res_s8196#*=}"
pass "s8196model.resourceKey=${MODEL_RES}"

# Must not broaden these ignores to all of internal/format or the whole tree.
if grep -E '^sonar\.issue\.ignore\.multicriteria\.s1186model\.resourceKey=' "${PROPS}" | grep -Eq '\*\*/internal/format/\*\*|^\*=\*\*$|\*\*/\*\.go'; then
  fail "s1186model.resourceKey must stay **/internal/format/model.go only"
fi
if grep -E '^sonar\.issue\.ignore\.multicriteria\.s8196model\.resourceKey=' "${PROPS}" | grep -Eq '\*\*/internal/format/\*\*|^\*=\*\*$|\*\*/\*\.go'; then
  fail "s8196model.resourceKey must stay **/internal/format/model.go only"
fi
pass "S1186/S8196 ignores scoped to model.go only"

# Why-comments must be present (not silent suppresses).
grep -q 'E4-S06' "${PROPS}" || fail "sonar-project.properties missing E4-S06 why-comment"
grep -qi 'marker' "${PROPS}" || fail "sonar-project.properties missing marker-method rationale"
grep -Eqi 'Block|Inline|-er|seal' "${PROPS}" || fail "sonar-project.properties missing S8196 rename rationale"
pass "E4-S06 why-comments present"

echo "All sonar_e4_s06 scope tests passed."
