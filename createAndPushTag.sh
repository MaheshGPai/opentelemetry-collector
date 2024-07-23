#!/bin/bash
#
# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

TAG=$1
unset GREP_OPTIONS
unset GREP_COLOR
if [ "${TAG}" == "" ]
then
  echo
  echo Unable to retrive the version
  echo
  exit 1
fi

function createAndPushTag() {
    git tag -d ${TAG}
    git push --delete intuit_remote ${TAG}

    git tag ${TAG}
    git push intuit_remote ${TAG}

    git tag -d config/confighttp/${TAG}
    git push --delete intuit_remote config/confighttp/${TAG}

    git tag config/confighttp/${TAG}
    git push intuit_remote config/confighttp/${TAG}

    git tag -d receiver/otlpreceiver/${TAG}
    git push --delete intuit_remote receiver/otlpreceiver/${TAG}

    git tag receiver/otlpreceiver/${TAG}
    git push intuit_remote receiver/otlpreceiver/${TAG}

    git tag -d pdata/${TAG}
    git push --delete intuit_remote pdata/${TAG}

    git tag pdata/${TAG}
    git push intuit_remote pdata/${TAG}
}

createAndPushTag
