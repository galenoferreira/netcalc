#!/usr/bin/env bash
set -euo pipefail

# Remember current directory
CUR_DIR=$(pwd)

# Change to repository directory (ensure REP_DIR is set or use current)
cd "${REP_DIR:-$CUR_DIR}" || exit 1

# Ensure .giti_tag exists with initial tag v0.0.0
if [[ ! -f .giti_tag ]]; then
  echo "v0.0.0" > .giti_tag
fi

# Read the last tag from .giti_tag
LAST_TAG=$(tail -n1 .giti_tag)

# Strip leading 'v' and split into components
VER=${LAST_TAG#v}
IFS='.' read -r MAJOR MINOR PATCH <<< "$VER"

# Increment patch version for the new tag
PATCH=$((PATCH + 1))
NEW_TAG="v${MAJOR}.${MINOR}.${PATCH}"

# Record the new tag in .giti_tag (append to preserve history)
echo "$NEW_TAG" >> .giti_tag

# Prepare prefix for commit message
PREFIX="[$NEW_TAG]"

# Show diff
echo
git diff

# Read commit message from user
# shellcheck disable=SC2162
read -p "Commit Message: " MSG

# Prefix the commit message
COMMIT_MSG="$PREFIX $MSG"

# Stage all changes
./build
git add .

# Create GitHub release
gh release create "$NEW_TAG" \
  --title "$COMMIT_MSG"

# Ensure we are on the main branch
git checkout main

# Create and push the new tag on main branch
git tag "$NEW_TAG"
git commit -m "$COMMIT_MSG"
git push origin main

# Upload each binary file in the bin subdirectories
for artifact in bin/*/*; do
  if [ -f "$artifact" ]; then
    gh release upload "$NEW_TAG" "$artifact" --clobber
  fi
done

git add .

# Return to original directory
cd "$CUR_DIR" || exit 1
