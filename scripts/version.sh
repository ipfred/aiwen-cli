#!/usr/bin/env bash
#
# version.sh — bump version in package.json, commit, and create git tag.
# package.json is the single source of truth for the version.
#
# Usage:
#   ./scripts/version.sh patch       # 0.1.0 -> 0.1.1
#   ./scripts/version.sh minor       # 0.1.0 -> 0.2.0
#   ./scripts/version.sh major       # 1.0.0 -> 2.0.0
#   ./scripts/version.sh 1.2.3      # set explicit version
#
# This script:
#   1. Reads current version from package.json
#   2. Bumps according to semver rule
#   3. Updates package.json
#   4. Commits and creates git tag v<version>
#   5. Pushes commit and tag to remote
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
PACKAGE_JSON="$PROJECT_DIR/package.json"

current_version() {
  node -p "require('$PACKAGE_JSON').version"
}

bump_version() {
  local current="$1"
  local rule="$2"
  local major minor patch

  IFS='.' read -r major minor patch <<< "$current"

  case "$rule" in
    patch)
      patch=$((patch + 1))
      ;;
    minor)
      minor=$((minor + 1))
      patch=0
      ;;
    major)
      major=$((major + 1))
      minor=0
      patch=0
      ;;
    *)
      # treat as explicit version string
      echo "$rule"
      return
      ;;
  esac

  echo "${major}.${minor}.${patch}"
}

if [ $# -lt 1 ]; then
  echo "Usage: $0 <patch|minor|major|version>"
  echo ""
  echo "  patch   Bump patch version (0.1.0 -> 0.1.1)"
  echo "  minor   Bump minor version (0.1.0 -> 0.2.0)"
  echo "  major   Bump major version (1.0.0 -> 2.0.0)"
  echo "  1.2.3   Set explicit version"
  echo ""
  echo "Current version: $(current_version)"
  exit 1
fi

RULE="$1"
CURRENT="$(current_version)"
NEW="$(bump_version "$CURRENT" "$RULE")"

if [ "$CURRENT" = "$NEW" ]; then
  echo "Version unchanged: $CURRENT"
  exit 0
fi

echo "Bumping version: $CURRENT -> $NEW"

cd "$PROJECT_DIR"

# Update package.json version using npm version (handles the commit + tag)
# If npm is not available, fall back to manual update
if command -v npm &>/dev/null; then
  npm version "$NEW" --no-git-tag-version 2>/dev/null || true
fi

# Verify package.json was updated
UPDATED="$(current_version)"
if [ "$UPDATED" != "$NEW" ]; then
  # Manual fallback — use node to update package.json
  node -e "
    const fs = require('fs');
    const pkg = require('$PACKAGE_JSON');
    pkg.version = '$NEW';
    fs.writeFileSync('$PACKAGE_JSON', JSON.stringify(pkg, null, 2) + '\n');
  "
fi

echo "Updated package.json to v$NEW"

# Check if package.json has uncommitted changes
if git diff --quiet package.json; then
  echo "package.json unchanged."
  exit 0
fi

# Commit
git add package.json
git commit -m "chore: bump version to v$NEW"

# Create tag
TAG="v$NEW"

# Check if tag already exists
if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Tag $TAG already exists locally, skipping."
else
  git tag "$TAG"
  echo "Created tag $TAG"
fi

echo ""
echo "Version updated to v$NEW"
echo ""
echo "Pushing commit and tag to remote..."
git push && git push origin "$TAG"

echo ""
echo "Done. CI will build and publish v$NEW."
