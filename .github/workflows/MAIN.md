## TL;DR; Quick Test

Run these commands to quickly test each part of the CI/CD pipeline:

```bash
# Ensure GitHub CLI is authenticated
gh auth login

# 1. Trigger Build & Test (push to main)
git checkout main
git pull origin main
git commit --allow-empty -m "ci: test build"
git push origin main

# 2. Trigger Build & Test on a Pull Request
git checkout -b test/ci
git push -u origin test/ci
gh pr create --title "Test CI" --body "Testing build & test" --base main

# 3. Trigger Changelog Update (publish a release)
git tag v0.0.0-test -m "Test Release"
git push origin v0.0.0-test
gh release create v0.0.0-test --notes "Testing changelog"

# 4. Trigger Deploy (after a successful build/test on main)
# The deploy job runs automatically on push to main after tests pass.
```

# PROD-CI/CD Pipeline Documentation

This document explains how to configure and use the consolidated CI/CD pipeline defined in `.github/workflows/main.yml`.

---

## 1. Workflow Name & Triggers

- **Name:** `CI / Changelog / Deploy`
- **Events:**
    - **push** to the `main` branch
    - **pull_request** targeting `main`
    - **release.published** (any new GitHub Release)

These triggers ensure:

1. Every change to `main` is built and tested.
2. Pull requests against `main` also run the build/tests.
3. Changelog is auto-updated when a Release is published.

---

## 2. Permissions & Defaults

```yaml
permissions:
  contents: write
  pull-requests: write

defaults:
  run:
    shell: bash
```

- **contents: write** — required to commit updated files (e.g., `CHANGELOG.md`).
- **pull-requests: write** — reserved for potential PR updates.
- All `run` steps execute under the `bash` shell by default.

---

## 3. Global Environment Variables

Defined in the `env:` section of the workflow:

```yaml
env:
  GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  COMMIT_USER: "github-actions[bot]"
  COMMIT_EMAIL: "actions@github.com"
```

- **GH_TOKEN:** built-in Actions token for API calls.
- **COMMIT_USER / COMMIT_EMAIL:** identity used for automated commits.

---

## 4. Job: Build & Test

```yaml
jobs:
  build_and_test:
    if: github.event_name != 'release'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/cache@v3  # cache Go modules
      - uses: actions/setup-go@v4
      - run: go build -v ./...
      - run: go test -v ./...
```

- **When it runs:** on `push` or `pull_request`, not on `release`.
- **Caching:** uses Go module cache to speed up builds.
- **Concurrency:** configured to cancel in-progress runs for the same branch.

---

## 5. Job: Update CHANGELOG

```yaml
jobs:
  update_changelog:
    if: github.event_name == 'release'
    runs-on: ubuntu-22.04
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Configure Git
        run: |
          git config user.name "${{ env.COMMIT_USER }}"
          git config user.email "${{ env.COMMIT_EMAIL }}"
          git checkout main
      - run: sudo apt-get update && sudo apt-get install -y gh
      - name: Generate CHANGELOG.md
        run: |
          gh release view "${{ github.event.release.tag_name }}" --json body --jq .body > CHANGELOG.md
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      - name: Commit & Push CHANGELOG
        run: |
          git add CHANGELOG.md
          if ! git diff --staged --quiet; then
            git commit -m "docs: update CHANGELOG for ${{ github.event.release.tag_name }}"
            git push origin main
          fi
```

- **When it runs:** only on `release.published`.
- **fetch-depth: 0:** ensures full history for tagging.
- **gh CLI:** used to extract release body automatically.

---

## 6. Job: Deploy to Production

```yaml
jobs:
  deploy:
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    needs: build_and_test
    steps:
      - uses: actions/checkout@v4
      - run: |
          chmod +x ./deploy.py
          ./deploy.py --build --test --commit "CI: deploy to production" --tag "deploy-${{ github.run_number }}"
        env:
          GITHUB_REPO: ${{ github.repository }}
          REP_DIR: ${{ github.workspace }}
          GOPATH: ${{ env.GOPATH }}
```

- **When it runs:** on every push to `main` after a successful build/test.
- **Custom script:** calls your `deploy.py` with environment variables.

---

## 7. Customization & Best Practices

1. **Branch protection**: require status checks before merging to `main`.
2. **Secrets management**: add additional secrets under repository Settings → Secrets.
3. **Reusable workflows**: extract common steps into separate `.github/workflows/*.yml` templates.
4. **Advanced caching**: cache Docker layers, npm modules, or other dependencies.
5. **Notifications**: integrate with Slack, email, or other channels for job failures.

This documentation will help users understand, maintain, and extend the consolidated CI/CD pipeline in `main.yml`.