#!/bin/bash

downloadPath=$1
urlFilePath=$2

mkdir -p ${downloadPath}
cd ${downloadPath}

# 由于下面两个过程都需要按行获取，将IFS修改成换行符可以避免空格将一行拆分
IFS=$'\n'
for url in `cat ${urlFilePath}`
do
  you-get ${url}
done

for file in `ls "${downloadPath}"`
do
  rclone copy -v "${downloadPath}/${file}" "one:${ONEDRIVE_BASE_PATH}"
done
