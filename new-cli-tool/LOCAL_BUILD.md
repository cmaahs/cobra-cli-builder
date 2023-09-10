# Local Build

Use this set of commands to perform a local build for tesing.

```bash
SEMVER=v0.0.999; echo ${SEMVER}
BUILD_DATE=$(gdate --utc +%FT%T.%3NZ); echo ${BUILD_DATE}
# GIT_COMMIT=$(git rev-parse HEAD); echo ${GIT_COMMIT}

PROGRAM=new-cli-tool
# go build -ldflags "-X ${PROGRAM}/main.semVer=${SEMVER} -X ${PROGRAM}/main.buildDate=${BUILD_DATE} -X ${PROGRAM}/main.gitCommit=${GIT_COMMIT} -X ${PROGRAM}/main.gitRef=/refs/tags/${SEMVER}"
go build -ldflags "-X ${PROGRAM}/version.semVer=${SEMVER} -X ${PROGRAM}/version.buildDate=${BUILD_DATE}"
if [[ -d ~/tbin ]]; then
  cp ${PROGRAM} ~/tbin/
fi
```
