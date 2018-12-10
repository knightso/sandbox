#! /bin/sh

BASEDIR=$(cd $(dirname $(dirname "$0")) && pwd)
GO111MODULE=on

gcloud beta emulators datastore start --no-store-on-disk --consistency=1.0 --quiet &
EMU_PID=$!; sleep 5; $(gcloud beta emulators datastore env-init)
trap 'pgrep -P $EMU_PID | xargs pgrep -P | xargs kill' 0 1 2 3 15

EMU_PORT=`echo $DATASTORE_EMULATOR_HOST | sed -e 's/.*://'`
cd $BASEDIR && dev_appserver.py --support_datastore_emulator=true --datastore_emulator_port=$EMU_PORT --env_var GOOGLE_CLOUD_PROJECT=$DATASTORE_PROJECT_ID ./

