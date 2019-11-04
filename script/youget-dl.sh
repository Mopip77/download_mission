#/bin/bash

timestamp=$1
urlFilepath=$2

urls=(`cat $2`)
downloadPath=${HOME}/${timestamp}
mkdir -p ${downloadPath}
cd ${downloadPath}

for url in ${urls[@]}
do
  you-get ${url}
done

for file in "`ls ${downloadPath}`"
do
  rclone copy "${file}" "one:${ONEDRIVE_BASE_PATH}"
done