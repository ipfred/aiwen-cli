#!/usr/bin/env bash
#
# version.sh — bump version and create git tag, similar to npm version
#
# Usage:
#   ./scripts/version.sh patch   # 0.1.0 -> 0.1.1
#   ./scripts/version.sh minor   # 0.1.0 -> 0.2.0
#   ./scripts/version.sh major   # 1.0.0 -> 2.0.0
#   ./scripts/version.sh 1.2.3  # set specific version
#
# This script:
#   1. Reads current version from VERSION file
#   2. Bumps according to semver rule
#   3. Updates VERSION file
#   4. Commits and creates git tag v<version>
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION_FILE="$PROJECT_DIR/VERSION"

current_version() {
  cat "$VERSION_FILE"
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

echo "Bumping version: $CURRENT -> $NEW"
echo "$NEW" > "$VERSION_FILE"

cd "$PROJECT_DIR"

if git diff --quiet VERSION; then
  echo "Version unchanged."
  exit 0
fi

git add VERSION
git commit -m "chore: bump version to v$NEW"
git tag "v$NEW"

echo ""
echo "Version updated to v$NEW"
echo ""
echo "Next steps:"
echo "  git push && git push --tags          # push commit and tag"
echo "  # or push tag to trigger CI release:"
echo "  git push origin v$NEW"
