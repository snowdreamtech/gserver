#!/bin/sh

# environment
export ProjectName="go_http_server"
Author="Yang Hui <sn0wdr1am@qq.com>"
LICENSE="Copyright (c) 2023-present SnowdreamTech Inc."
CmdPath="github.com/snowdreamtech/gserver/pkg/env"
BuildTime=$(date +'%Y.%m.%d.%H%M%S%z')
CommitHash=N/A
CommitHashFull=N/A
GoVersion=N/A
GitTag=N/A
Debug=false

if [[ $(go version) =~ [0-9]+\.[0-9]+\.[0-9]+ ]];
then
    GoVersion=${BASH_REMATCH[0]}
fi

GV=$(git tag -l | sort -V --reverse || echo 'N/A')
if [[ $GV =~ [^[:space:]]+ ]];
then
    GitTag=${BASH_REMATCH[0]}
fi

GH=$(git log -1 --pretty=format:%h || echo 'N/A')
if [[ GH =~ 'fatal' ]];
then
    CommitHash=N/A
else
    CommitHash=$GH
fi

GHF=$(git log -1 --pretty=format:%H || echo 'N/A')
if [[ GHF =~ 'fatal' ]];
then
    CommitHashFull=N/A
else
    CommitHashFull=$GHF
fi

LDGOFLAGS="-X '$CmdPath.ProjectName=$ProjectName'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.Author=$Author'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.BuildTime=$BuildTime'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.CommitHash=$CommitHash'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.CommitHashFull=$CommitHashFull'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.GoVersion=$GoVersion'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.GitTag=$GitTag'"
LDGOFLAGS="$LDGOFLAGS -X '$CmdPath.LICENSE=$LICENSE'"

if [ "$Debug" = false ] ; then
    LDGOFLAGS="$LDGOFLAGS  -s -w"
fi

# echo $LDGOFLAGS
export LDGOFLAGS
