# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
Release notes are generated from [Conventional Commits](https://www.conventionalcommits.org/)
on the default branch using [git-cliff](https://git-cliff.org/).

## [Unreleased]

### Bug Fixes

- **sonar:** Exclude tests from sources so CPD clears [094fc8e](https://github.com/platformrelay/kollect-render/commit/094fc8ee9a97f1647fee02d29d3036b25930cd4a)

- **ci:** Disable Codecov coverage file search [1107268](https://github.com/platformrelay/kollect-render/commit/110726825c04ae240bc2585b26db5c387ec61885)

- **release:** Match CHANGELOG VERSION as fixed string [ba9cdf5](https://github.com/platformrelay/kollect-render/commit/ba9cdf5f9adc2286d165251ae3d91cfe31ac4b73)

- **ci:** Scope docs Pages OIDC perms to deploy job [b87637b](https://github.com/platformrelay/kollect-render/commit/b87637b630f0e6942a1ca4511613a4118baf875c)

- **hack:** Enforce HTTPS on git-cliff curl download [5035c23](https://github.com/platformrelay/kollect-render/commit/5035c23e1fffe680f9b76f43c067cb8c75d9dd9c)

- **docs:** Correct template threat model in SECURITY [98df359](https://github.com/platformrelay/kollect-render/commit/98df359493e3b5b6c453bfaacb8e06743b85330b)

- **ci:** SHA-pin workflow actions and HTTPS-harden gitleaks fetch [134399b](https://github.com/platformrelay/kollect-render/commit/134399bf8275be0acc308b3ee9bf2264c5e44b84)


### Refactoring

- **format:** Split writeConfluenceBlock under cognitive 15 ([#10](https://github.com/platformrelay/kollect-render/pull/10))[d8a3cfb](https://github.com/platformrelay/kollect-render/commit/d8a3cfbaff56634ff8e4d59f3e944349e0358d09)

- **format:** Split EnvInventoryModel under gocyclo 20 ([#7](https://github.com/platformrelay/kollect-render/pull/7))[4fc0405](https://github.com/platformrelay/kollect-render/commit/4fc04056e601bb0aa704cfe0fb337cd9018289d0)

- **cli:** Split runRender under gocyclo 20 ([#4](https://github.com/platformrelay/kollect-render/pull/4))[9428c6c](https://github.com/platformrelay/kollect-render/commit/9428c6c0b5737fe1c1b0bc4a8f5657b039bbf657)

## [0.1.0] - 2026-07-25

### Bug Fixes

- Reject non-markdown format with --template [fa5461f](https://github.com/platformrelay/kollect-render/commit/fa5461f7c727384e578bc10f790f66af70c78b6d)

- Expand REWE-trace forbidden-marker coverage [9818524](https://github.com/platformrelay/kollect-render/commit/9818524d876decb9a8751587009aac1239cce503)


### Features

- Publish multi-platform binaries and GHCR image on tags [aab2abb](https://github.com/platformrelay/kollect-render/commit/aab2abbe0ed56c57ea34e47c84d9992a6b189825)

- Enforce kollect-parity quality gates at 90% coverage [53004ab](https://github.com/platformrelay/kollect-render/commit/53004abc1fb788b5ff567aca10731c86613d2f2e)

- Wire render CLI override flags [53a9950](https://github.com/platformrelay/kollect-render/commit/53a995018786cf83968c110ca1e0b0fc944d7f88)

- Emit digest/metadata render sidecar [aca5e68](https://github.com/platformrelay/kollect-render/commit/aca5e683e2c1874e9cf5c1f940d4e18f506c6808)

- Implement format registry and encoders [f322fbc](https://github.com/platformrelay/kollect-render/commit/f322fbc10e9c5e641f72a6f9b61c0d36b5391080)

- Implement template engine and helper set [bf21955](https://github.com/platformrelay/kollect-render/commit/bf21955212201743cae5d3e1f2f82e88a4d78d8b)

- Implement inventory schema validate command [9d2ae11](https://github.com/platformrelay/kollect-render/commit/9d2ae11cbc219458a63a88f1b8ee27237f7feba2)

