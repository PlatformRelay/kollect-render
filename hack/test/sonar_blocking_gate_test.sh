#!/usr/bin/env bash
# Assert the SonarCloud job is a BLOCKING gate (OD-3 revisit, 2026-08-06).
#
# OD-3 made `sonarqube` non-blocking so the quality gate could not veto the very
# lanes (E4-S01/E4-S02) that were fixing it. That gate went OK once Automatic
# Analysis was disabled and CI analysis actually ran, so the escape hatch came
# out. Blocking here means BOTH halves, because either alone is a false sense of
# safety:
#
#   1. no `continue-on-error` — a scanner failure fails the job. Without this a
#      hard exit 3 ("running CI analysis while Automatic Analysis is enabled")
#      was swallowed for a day while the workflow reported success.
#   2. `sonar.qualitygate.wait=true` — the scanner polls for the gate verdict
#      instead of fire-and-forget. Without this the job goes green the moment the
#      report is uploaded, so a gate that flips to ERROR is never noticed.
#
# Run: bash hack/test/sonar_blocking_gate_test.sh
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
wf="${root}/.github/workflows/ci.yaml"
props="${root}/sonar-project.properties"

fail=0

# 1. No job or step in ci.yaml may opt out of failing the workflow. Deliberately
#    whole-file: nothing in this workflow has a legitimate use for it today.
#    Matches bare, quoted and expression forms.
if grep -Eq "^[[:space:]]*continue-on-error:[[:space:]]*(\"?true\"?|'true'|\\\$\\{\\{)" "$wf"; then
  echo "FAIL: continue-on-error opt-out still present in ci.yaml" >&2
  grep -nE '^[[:space:]]*continue-on-error:' "$wf" >&2
  fail=1
else
  echo "ok: no continue-on-error opt-out in ci.yaml"
fi

# 2. The sonarqube job must exist and keep needing `check` — the scan consumes
#    coverage.out via the artifact that `check` uploads.
if grep -Eq '^[[:space:]]{2}sonarqube:' "$wf"; then
  echo "ok: sonarqube job present"
else
  echo "FAIL: sonarqube job missing from ci.yaml" >&2
  fail=1
fi

if grep -Eq '^[[:space:]]*needs:[[:space:]]*\[[[:space:]]*check[[:space:]]*\]' "$wf"; then
  echo "ok: sonarqube needs [check] (coverage artifact producer)"
else
  echo "FAIL: sonarqube no longer declares needs: [check]" >&2
  fail=1
fi

# 3. The scanner must wait for the quality gate verdict, not merely upload.
if grep -Eq '^sonar\.qualitygate\.wait[[:space:]]*=[[:space:]]*true' "$props"; then
  echo "ok: sonar.qualitygate.wait=true"
else
  echo "FAIL: sonar.qualitygate.wait=true missing from sonar-project.properties" >&2
  echo "      without it the job passes on upload and never sees an ERROR gate" >&2
  fail=1
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "sonar_blocking_gate_test: RED" >&2
  exit 1
fi
echo "sonar_blocking_gate_test: GREEN"
