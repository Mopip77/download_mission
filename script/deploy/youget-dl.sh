#/bin/bash

downloadPath=$1
urlFilepath=$2

urls=(`cat $2`)
mkdir -p ${downloadPath}
cd ${downloadPath}

for url in ${urls[@]}
do
  you-get ${url}
done

# 由于文件名可能有空格，所以使用一个基本不可能出现的字符替换后再替换回来
symbol="觉d怼e部z科k恁"
for file in `ls "${downloadPath}" | sed 's/ /'"${symbol}"'/g'`
do
  realFileName=`sed 's/'"${symbol}"'/ /g' <<<$file`
  rclone copy -v "${downloadPath}/${realFileName}" "one:${ONEDRIVE_BASE_PATH}"
done
