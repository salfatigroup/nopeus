#!/bin/sh
set -e

if test "$DISTRIBUTION" = "pro"; then
    echo "Installing nopeus pro 🎉"
    RELEASES_URL="https://github.com/salfatigroup/nopeus/releases"
    FILE_BASENAME="nopeus"
else
    echo "Only nopeus pro is supported at this time."
    exit 1
fi

test -z "$VERSION" && VERSION="$(curl -sfL -o /dev/null -w %{url_effective} -u elonsalfati:$NOPEUS_TOKEN "$RELEASES_URL/latest" |
    rev |
    cut -f1 -d'/'|
    rev)"

test -z "$VERSION" && {
    echo "Unable to get nopeus version." >&2
    exit 1
}

test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
export TAR_FILE="$TMPDIR/${FILE_BASENAME}_$(uname -s)_$(uname -m).tar.gz"

(
    cd "$TMPDIR"
    echo "Downloading nopeus $VERSION..."
    curl -sfLo -u elonsalfati:$NOPEUS_TOKEN "$TAR_FILE" \
        "$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_$(uname -s)_$(uname -m).tar.gz"
    curl -sfLo -u elonsalfati:$NOPEUS_TOKEN "checksums.txt" "$RELEASES_URL/download/$VERSION/checksums.txt"
    curl -sfLo -u elonsalfati:$NOPEUS_TOKEN "checksums.txt.sig" "$RELEASES_URL/download/$VERSION/checksums.txt.sig"
    echo "Verifying checksums..."
    sha256sum --ignore-missing --quiet --check checksums.txt
    if command -v cosign >/dev/null 2>&1; then
        echo "Verifying signatures..."
        COSIGN_EXPERIMENTAL=1 cosign verify-blob \
            --signature checksums.txt.sig \
            checksums.txt
    else
        echo "Could not verify signatures, cosign is not installed."
    fi
)
