
## CI/CD Pipeline Documentation

This document explains how to configure and use the consolidated CI/CD pipeline defined in `.github/workflows/main.yml`.

---

## 1. Workflow Name & Triggers

- **Name:** `MAIN CD-CI Workflow`
- **Events:**
    - **release.published** (when a new GitHub Release is published)

These triggers ensure:

1. Changelog is auto-updated when a Release is published.
2. Deployment to production runs on release tags.

---

## 2. Job: Update CHANGELOG

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
- **fetch-depth: 0:** ensures full history for proper tagging.
- **gh CLI:** automatically extracts the release body.

---

## 3. Job: Deploy to Production

```yaml
jobs:
  deploy:
    if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v')
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

- **When it runs:** on push of tags matching `v*`, after a successful build/test.
- **Custom script:** invokes `deploy.py` with necessary environment variables.

---

## 7. Customization & Best Practices

1. **Branch protection**: require status checks before merging to `main`.
2. **Secrets management**: store sensitive values under repository Settings → Secrets.
3. **Reusable workflows**: extract common steps into shared `.github/workflows/*.yml`.
4. **Advanced caching**: leverage caches for Docker layers, npm modules, or other dependencies.
5. **Notifications**: integrate with external services (Slack, email) for failure alerts.

This documentation will help users understand, maintain, and extend the consolidated CI/CD pipeline in `main.yml`.