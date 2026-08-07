# Demo — AI-assisted Renovate PRs

This branch is a **walkthrough only** (nothing here runs in CI). It shows, with
a worked example, what the two workflows on the `dep-bumping-ai` branch do when
Renovate opens a dependency-bump PR.

The two workflows:

| Workflow | Fires when | Does |
|---|---|---|
| `renovate-changelog-summary.yaml` | Renovate opens/updates a PR | Posts a Changelog Risk Summary comment; labels `breaking-change` if warranted. Read-only. |
| `renovate-ci-fixer.yaml` | `Lint` or `Build Examples` **fails** on a Renovate PR | Makes **one** attempt to fix Go code, pushes to the PR branch, or escalates via comment. |

---

## Scenario

Renovate opens a PR bumping a Go dependency used in `tests/`:

> **⬆️ Update module github.com/example/widget to v2**

### Step 1 — Changelog Summary runs (every Renovate PR)

Claude (Haiku) reads the changelog Renovate embedded and posts one comment:

> ## Changelog Risk Summary
> - **Breaking changes:** `widget.New()` now requires a `context.Context` as its
>   first argument; the `widget.Legacy` type was removed.
> - **Security fixes:** None found.
> - **Deprecations:** `widget.SetTimeout` is deprecated in favour of context deadlines.
> - **Reviewer note:** Breaking API change — the `tests/` package will need a small
>   adaptation before this is mergeable.

Because a real breaking change was found, the PR also gets the **`breaking-change`** label.

### Step 2 — CI fails

`Build Examples` goes red because `tests/` no longer compiles:

```
tests/widget_test.go:14:22: not enough arguments in call to widget.New
	have ()
	want (context.Context)
```

See [`before/widget_test.go`](before/widget_test.go) — the pre-bump code.

### Step 3 — CI Fixer runs (only on that failure)

Claude (Sonnet) checks out the Renovate branch, reproduces the failure, and makes
the **minimal** fix — passing a context, within the allowed files only
(`tests/**/*.go`, `go.mod`, `go.sum`). See [`after/widget_test.go`](after/widget_test.go).

It pushes a signed-off commit to the PR branch:

```
fix: adapt to github.com/example/widget v2 bump

Signed-off-by: ...
```

and comments what it changed. CI re-runs green → the PR is ready for a human to merge.

### If it *couldn't* fix it

(e.g. a breaking major bump needing changes outside the allowed files) it edits
**nothing** and posts a single comment: what failed, what it tried, and that a
human needs to take over.

---

## Guardrails demonstrated here

- **Author-gated:** both only act on `renovate[bot]` PRs.
- **Failure-scoped:** the fixer only reacts to `Lint` / `Build Examples`; QEMU
  e2e, Scorecard, and OSV-Scanner are excluded so it never papers over flaky infra.
- **Path allowlist:** the fixer can touch only `tests/**/*.go`, `go.mod`, `go.sum`
  — never workflows, `renovate.json`, or secrets.
- **Single attempt:** one try per failure, then it stops and escalates.
