#!/bin/bash

# 进入脚本文件目录，因为之后的脚本都在同级目录
cd `dirname $0`
# install aria2
wget -N --no-check-certificate https://raw.githubusercontent.com/ToyoDAdoubi/doubi/master/aria2.sh && chmod +x aria2.sh
./aria.expect
# config aria2
mkdir -p /data/Download
sed -i 's@dir=.*@dir=/data/Download@' ${HOME}/.aria2/aria2.conf
sed -i 's@rpc-secret=.*@rpc-secret='"$ARIA_RPC_PWD"'@' ${HOME}/.aria2/aria2.conf
echo "on-download-complete=${HOME}/.aria2/upload2one.sh" >> ${HOME}/.aria2/aria2.conf
/etc/init.d/aria2 restart
rm aria2.sh

echo -e "#!/bin/bash
filepath=\$3	 #取文件原始路径，如果是单文件则为/Download/a.mp4，如果是文件夹则该值为文件夹内第一个文件比如/Download/a/1.mp4
path=\${3%/*}	 #取文件根路径，如把/Download/a/1.mp4变成/Download/a
downloadpath='/data/Download'	#Aria2下载目录
name='one' #配置Rclone时的name
folder='${ONEDRIVE_BASE_PATH}'	 #网盘里的文件夹，如果是根目录直接留空
MinSize='10k'	 #限制最低上传大小，默认10k，BT下载时可防止上传其他无用文件。会删除文件，谨慎设置。
MaxSize='15G'	 #限制最高文件大小，默认15G，OneDrive上传限制。

if [ \$2 -eq 0 ]; then exit 0; fi

while true; do
if [ \"\$path\" = \"\$downloadpath\" ] && [ \$2 -eq 1 ]	#如果下载的是单个文件
    then
    rclone move -v \"\$filepath\" \${name}:\${folder} --min-size \$MinSize --max-size \$MaxSize
    rm -vf \"\$filepath\".aria2	#删除残留的.aria.2文件
    rm \"\$filepath\"
    exit 0
elif [ \"\$path\" != \"\$downloadpath\" ]	#如果下载的是文件夹
    then
    while [[ \"\`ls -A \"\$path/\"\`\" != \"\" ]]; do
    rclone move -v \"\$path\" \${name}:/\${folder}/\"\${path##*/}\" --min-size \$MinSize --max-size \$MaxSize --delete-empty-src-dirs
    rclone delete -v \"\$path\" --max-size \$MinSize	#删除多余的文件
    rclone rmdirs -v \"\$downloadpath\" --leave-root	#删除空目录，--delete-empty-src-dirs 参数已实现，加上无所谓。
    done
    rm -vf \"\$path\".aria2	#删除残留的.aria2文件
    rm -rf \"\$path\"
    exit 0
fi
done" > ${HOME}/.aria2/upload2one.sh
chmod +x ${HOME}/.aria2/upload2one.sh

# install and config rclone
wget https://rclone.org/install.sh
bash install.sh
rm install.sh
./rclone.expect "${ONEDRIVE_TOKEN}"

# start download api
cd ${HOME}/api
./dl_api