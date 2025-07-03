## TL;DR; Quick Test

Run these commands to quickly test each part of the CI/CD pipeline:

```bash
gh auth login
```

```bash
git checkout main
```

```bash
git pull origin main
```

```bash
git commit --allow-empty -m "ci: test build"
```

```bash
git push origin main
```

```bash
git checkout -b test/ci
```

```bash
git push -u origin test/ci
```

```bash
gh pr create --title "Test CI" --body "Testing build & test" --base main
```

```bash
git tag v0.0.0-test -m "Test Release"
```

```bash
git push origin v0.0.0-test
```

```bash
gh release create v0.0.0-test --notes "Testing changelog"
```

```bash
# The deploy job runs automatically on push to main after tests pass.
```