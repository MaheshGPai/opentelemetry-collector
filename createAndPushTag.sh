#! /bin/bash

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
    git push --delete origin ${TAG}

    git tag ${TAG}
    git push origin ${TAG}

    git tag -d config/confighttp/${TAG}
    git push --delete origin config/confighttp/${TAG}

    git tag config/confighttp/${TAG}
    git push origin config/confighttp/${TAG}

    git tag -d receiver/otlpreceiver/${TAG}
    git push --delete origin receiver/otlpreceiver/${TAG}

    git tag receiver/otlpreceiver/${TAG}
    git push origin receiver/otlpreceiver/${TAG}

    git tag -d pdata/${TAG}
    git push --delete origin pdata/${TAG}

    git tag pdata/${TAG}
    git push origin pdata/${TAG}
}

createAndPushTag
