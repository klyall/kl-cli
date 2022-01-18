#! /bin/bash
git tag -a v$1 -m "Releasei $1"
git push origin v$1

goreleaser release --rm-dist

