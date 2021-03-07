#!/bin/sh
basepath=$(cd `dirname $0`; pwd)
echo "start..."
cd $basepath
nohup ./tgdown -d $basepath > /dev/null 2>&1 &
echo "The TG Download is running..."