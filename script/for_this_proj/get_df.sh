#!/bin/bash

# 先求使用比例，使用bc工具保留两位小数
res=(`df | awk '{if ($NF=="/") print $2,"\n",$3}'`)
echo `echo "scale=2; ${res[1]} / ${res[0]}" | bc`

# 求格式化后的使用量
df -h | awk '{if ($NF=="/") print $2,"\n",$3}'

