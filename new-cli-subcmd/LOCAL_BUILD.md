# Local Build

Use this set of commands to perform a local build for tesing.

```bash
SEMVER=v0.0.1; echo ${SEMVER}
BUILD_DATE=$(gdate --utc +%FT%T.%3NZ); echo ${BUILD_DATE}
# GIT_COMMIT=$(git rev-parse HEAD); echo ${GIT_COMMIT}

#go build -ldflags "-X new-cli-subcmd/cmd.semVer=${SEMVER} -X new-cli-subcmd/cmd.buildDate=${BUILD_DATE} -X new-cli-subcmd/cmd.gitCommit=${GIT_COMMIT} -X new-cli-subcmd/cmd.gitRef=/refs/tags/${SEMVER}" && \
go build -ldflags "-X new-cli-subcmd/cmd.semVer=${SEMVER} -X new-cli-subcmd/cmd.buildDate=${BUILD_DATE}" && \
./new-cli-subcmd version | jq .
```
