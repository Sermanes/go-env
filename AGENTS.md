# Agents — Operating Rules

Loaded before any agent acts on this repository. Applies to every session.

`go-env` is a small Go library: typed accessors over `os.LookupEnv` with default fallbacks. Public API lives in `env.go`. Tests live in `env_test.go`. There is no `cmd/`, no `internal/`, no Docker. Keep it that way unless explicitly asked.

## Hard rules

### TDD is non-negotiable
- Every feature, every bugfix: write the failing test **first**. No production code lands without a red test that becomes green.
- Order is always: red → green → refactor.
- If a bug is reported, the first commit (or the first hunk in the same commit) must add a test that reproduces the bug — proving the fix actually fixes something.
- Refactors without behavior change are the only exception, and even then existing tests must stay green throughout.

### Table-driven tests
- All tests use Go's table-driven pattern. Each scenario is a row in a `[]struct{...}` slice with named fields.
- Each row gets a unique `name` and is run via `t.Run(tc.name, func(t *testing.T) {...})`.
- Subtests use snake-case-with-spaces sentence names: `"when the key exists"`, `"when the value overflows uint8"`. Match the existing style in `env_test.go`.
- Use `t.Setenv` (never `os.Setenv`) so cleanup is automatic.
- Use `testify/assert` for assertions. Use `assert.Equal`, `assert.Nil`, etc. — match existing style.
- One behavior per row. Never combine "valid input AND invalid input" in the same row.
- Cover at minimum: happy path, missing variable, invalid value, edge cases (empty string, negative, overflow, whitespace).

### Go best practices
- `gofmt` + `goimports` on every file. Enforced by hooks and CI.
- `go vet` clean. Enforced by hooks and CI.
- `golangci-lint` (config `.golangci.yml`) clean. Enforced by hooks and CI.
- Exported identifier → docstring required. Form: `// Name does X. Returns Y when Z.` First word is the identifier name.
- Package comment required (`// Package env ...`).
- No magic numbers in production code — use named constants (`bitSize8`, `decimalBase`, etc.). Test files exempt.
- Errors: prefer `strconv.ParseUint`/`ParseInt`/`ParseFloat` with explicit `bitSize` over `Atoi`+cast. Bit-size casts must be range-checked or use the parser's bit-size argument.
- No silent overflow. `uint8` parsing must reject `>255` and negatives — `ParseUint(v, 10, 8)`, not `Atoi` + cast.
- Public API stability: never break existing function signatures. Add new functions; do not modify or remove existing ones without explicit instruction.
- No unnecessary abstractions. This is a leaf library. No interfaces unless there's a real second implementation.

### README is part of every change
- Any change that adds, removes, or modifies a public function **must** update `Readme.md` in the same commit.
- New function → add a `#### <Type>` section under `## Usage` with a code snippet matching the existing style.
- Removed function → delete its section.
- Behavior change visible to users (new validation, new error path) → update the function's snippet description.
- The `## Features` bullet list must stay accurate. Add a sub-bullet when relevant.
- A PR that touches `env.go` and not `Readme.md` is incomplete and the reviewer rejects it.

### Conventional Commits
- Format: `<type>(<optional scope>): <imperative subject>`. Subject ≤ 72 chars.
- Allowed types: `feat`, `fix`, `chore`, `refactor`, `test`, `docs`, `ci`, `build`, `perf`.
- Body explains *why*, not *what*. Bullet points OK.
- Breaking change → `!` after type (e.g. `feat!:`) and a `BREAKING CHANGE:` footer.
- One logical change per commit. No `WIP`, no `update`, no `fix stuff`.

### Verification gate
- `make ci` is the gate. It must pass locally before any commit, and on CI before any merge.
- `make ci` runs: `tidy → fmt → vet → test → lint`.
- `make before-push` adds the race detector.
- Hooks (`scripts/hooks/pre-commit`, `scripts/hooks/pre-push`) cannot be bypassed. Failure means fix the issue, never `--no-verify`.

### Files agents must not modify casually
Only touch when the user explicitly asks:
- `Makefile`
- `.golangci.yml`
- `.github/workflows/*`
- `scripts/hooks/*`
- `.pre-commit-config.yaml`
- `go.mod` / `go.sum` (except via `go mod tidy` or version bumps the user requested)

## Workflow per task

1. **Understand.** Read `env.go`, `env_test.go`, `Readme.md`. Confirm scope.
2. **Red test.** Add a table-driven row (or a new test function with a table) covering the new behavior. Run `make test` — confirm it fails for the expected reason.
3. **Green code.** Minimum implementation in `env.go` to make the row pass. Run `make test` — confirm green.
4. **Refactor.** Extract constants, simplify, deduplicate. Tests stay green.
5. **README.** Update `Readme.md` — add/modify usage block, update feature bullet list.
6. **Verify.** Run `make ci`. Must be clean (0 issues).
7. **Commit.** Conventional Commit message. Body explains why.

## Escalation

Stop and ask the human if any of these happen:
- An existing test fails on the clean baseline before your change.
- A hook fails and the cause is not obvious from the output.
- The change requires breaking the public API.
- `make ci` flags an issue you cannot fix without touching one of the protected files.
- README structure does not have a clear place for the new section.
