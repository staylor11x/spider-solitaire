# Log of Key Decisions in the Project

### Testing Predicimant [18/09/2025]

We are facing issues with testing internal/external methods, most of these issues seem to be stemming from the test tools library that we implemented. When we implemented this it caused us to have circular dependencies in our tests, we fixed this by moving the tests into a "different" package. 

The problem that this caused is that now we can only test "exposed" methods that are part of the public facing API, originally this was fine but as we started to introduce more complicated logic into the un-exported helper functions this issue quickly surfaced.

I am making the decision based on this to scrap the external test tools module and either simply include the helper functions for the tests within the test file they are required in or not use them at all.

I think this idea of test tool helper functions works better in Java esq languages

---

### Dual Versioning Strategy [24/01/2026]

Implemented a two-track versioning system:
- **Main branch:** Automatic patch version increments (v1.2.0 → v1.2.1 → v1.2.2)
- **Release branches:** Manual minor/major releases (release/v1.3.0 → v1.3.0 tag)

Rationale: Enables rapid development with tracked patches while maintaining controlled production releases. Developers know main always has a stable patchable version, users know releases are curated.

### Workflow Chaining with workflow_run [24/01/2026]

Chained release workflows using GitHub Actions' `workflow_run` trigger:
1. Release Branch Tagging → creates tag
2. Build Release Binaries (triggered by tag from step 1)
3. Create GitHub Release (triggered by build completion)

Rationale: Direct GITHUB_TOKEN tag pushes don't trigger workflows (security feature). Using workflow_run keeps builds automatic while respecting GitHub's security model. Alternative would be personal access tokens (less secure).

### No Code Signing Implementation [24/01/2026]

Decided against code signing certificates ($200-500/year) as first-pass implementation.

Rationale: Checksums provide integrity verification (what users really need). Windows "unverified publisher" warning is acceptable for open-source. Can add code signing later if distribution expands or sponsor emerges. Focus now: shipping working software.

### Linux + Windows Multi-Platform Support [24/01/2026]

Building on native runners (ubuntu-latest for Linux, windows-latest for Windows) rather than cross-compilation.

Rationale: Avoids CGO cross-compilation complexity with Ebiten. Native builds are faster, more reliable. Covers 95% of users (amd64). Can expand to arm64 later if needed. 