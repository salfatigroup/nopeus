#!/bin/bash -e

TOKEN="$NOPEUS_TOKEN"
OWNER="salfatigroup"
REPO="nopeus"
VERSION="${VERSION:-latest}"
CURRENT_LOCATION="$(pwd)"
FILE="nopeus_$(uname -s)_$(uname -m).tar.gz"

function gh_curl() {
  API_GITHUB="https://api.github.com"
  curl ${TOKEN:+-H "Authorization: token $TOKEN"} \
       -H "Accept: application/vnd.github.v3.raw" \
       -s $API_GITHUB$@
}

function gh_get_file() {
  file="$1"
  version="$2"
  asset_id=`gh_curl /repos/$OWNER/$REPO/releases/tags/$version | jq ".assets | map(select(.name == \"$file\"))[0].id"`
  if [ "$asset_id" = "null" ]; then
    >&2 echo -e "\033[31mERROR: \"$file\" not found at \"$version\"\033[0m"
    return 1
  fi

  wget -q --show-progress --auth-no-challenge --header='Accept:application/octet-stream' \
    https://${TOKEN:+$TOKEN:@}api.github.com/repos/$OWNER/$REPO/releases/assets/$asset_id \
    -O $file

  return $?
}

# check if required tools exists
if ! [ -x "$(command -v jq)" ]; then
  >&2 echo -e "\033[31mERROR: \"jq\" not found\033[0m"
  exit 1
fi

if ! [ -x "$(command -v wget)" ]; then
  >&2 echo -e "\033[31mERROR: \"wget\" not found\033[0m"
  exit 1
fi

if [ -z "$TOKEN" ]; then
  >&2 echo -e "\033[33mWARNING: \$TOKEN is empty\033[0m"
  exit 1
fi

if [ "$VERSION" = "latest" ]; then
  VERSION=`gh_curl /repos/$OWNER/$REPO/releases/latest | jq --raw-output ".tag_name"`
  echo "Latest version is \"$VERSION\""
fi

mkdir -p "/tmp/$REPO-$VERSION"
cd "/tmp/$REPO-$VERSION" > /dev/null

# Start by getting the checksums if they are available.
if ! gh_get_file checksums.txt $VERSION; then
  >&2 echo -e "\033[33mWARNING: Checksums will NOT be computed\033[0m"
fi

if ! [ -z "$FILE" ]; then
  if ! gh_get_file $FILE $VERSION; then
    exit 1
  fi

  # Verify the sha256 sum for this file only.
  grep "$FILE" checksums.txt | sha256sum --check
else
  release_id=`gh_curl /repos/$OWNER/$REPO/releases/tags/$VERSION | jq ".id"`
  assets=`gh_curl /repos/$OWNER/$REPO/releases/$release_id/assets | jq --raw-output '.[].name'`
  for asset in $assets; do
    if ! [ "$asset" = "checksums.txt" ]; then
      if ! gh_get_file $asset $VERSION; then
        exit 1
      fi
    fi
  done

  # Verify the sha256 sum for all files.
  sha256sum --check checksums.txt
fi

mkdir -p "$HOME/nopeus"
tar -xf "$FILE" -C "$HOME/nopeus/"
# install in /usr/local/bin
tar -xf "$FILE" -C /usr/local/bin/
echo "nopeus binary is located under $HOME/nopeus/nopeus"

cd $CURRENT_LOCATION > /dev/null
