## Setup
The repository's tools module owns the exact Verify version:
```sh
go -C tools get -tool github.com/secengcommons/verify/cmd/secverify@VERSION
```

## Use
```yaml
jobs:
  verify:
    uses: secengcommons/ci/.github/workflows/verify.yml@COMMIT

  linux:
    uses: secengcommons/ci/.github/workflows/linux.yml@COMMIT
```

## Checks
Workflow | Evidence
--- | ---
verify.yml | Static analysis, Go compatibility, complete coverage, race and fuzzing
linux.yml | Ubuntu AMD64 and ARM64 execution
windows.yml | Windows Server AMD64 tests and Windows 11 ARM64 inventory
macos.yml | macOS 15 and 26 on ARM64 and Intel
fuzz.yml | Linux, Windows and macOS fuzz campaigns
ports.yml | Portable root-module packages across every target reported by the selected Go toolchain
containers.yml | Root-module tests in a restricted Alpine Linux container
alpine.yml | Repository tests in an Alpine Linux virtual machine
bsd.yml | Root-module tests on FreeBSD, OpenBSD and NetBSD
illumos.yml | Root-module tests on OmniOS
solaris.yml | Root-module tests on Oracle Solaris
dependencies.yml | Pull-request dependency review
codeql.yml | Root-module Go and GitHub Actions security and quality queries
qualification.yml | Complete manual candidate run
