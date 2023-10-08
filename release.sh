TAGS=`git show-ref --tags | wc -l`;
RELEASE_NUMBER=$(expr $TAGS + 1);
echo release number: ${RELEASE_NUMBER};
VERSION=v0.0.${RELEASE_NUMBER}-alpha;
echo version tag: ${VERSION};

git tag ${VERSION}
git push origin ${VERSION}

goreleaser release --clean