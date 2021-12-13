#build bash file
#!/bin/bash
project=store_server_http
build_dir=/tmp/build_${project}
gopath=/data/workspace/${pipeline.id}
# 进入项目代码目录
cd ${WORKSPACE}/store_server_http

# 编译出二进制文件
if [[ ! -d ${build_dir} ]]; then
    mkdir ${build_dir}
else
    rm -rf ${build_dir}/*
fi
cd ${WORKSPACE}/store_server/cmd/$project
export GOPATH=${gopath}:/data/erichli/projects/
git fetch --depth=500 && gopack -n $project -u $build_dir

#复制目标构件到工作空间
cp ${build_dir}/*.rpm ${WORKSPACE}/store_server/

#*********************************************************#
#deploy bash file
#!/bin/bash
project=store_server_http
tmp_dir=/tmp/${project}
deploy_dir=/tmp/deploy_${project}
work_dir=/data/apps/store_server_http

#step 1 进入临时文件目录，解包
#cd ${tmp_dir}
#tar_file=${tmp_dir}/store_server.tgz
#tar -zxf $tar_file --strip-components 1 -C $dst_dir/HEAD
#if [[ ! -d $tmp_dir/src/store_server ]]; then
#    mkdir -p $tmp_dir/src/store_server
#fi
#tar -zxf $tar_file -C $tmp_dir/src/store_server

#step 2 拷贝app配置文件
if [[ ! -d ${work_dir}/configs ]]; then
    mkdir -p ${work_dir}/configs
fi    
cp -r /data/erichli/store_server_configs/configs/${project}_* ${work_dir}/configs/

#step 3 创建日志目录
if [[ ! -d ${work_dir}/logs ]]; then
    mkdir ${work_dir}/logs
fi

#step 4 link config by env
if [[ -L ${work_dir}/conf/${project}.yaml ]]; then
    unlink ${work_dir}/conf/${project}.yaml
fi

if [[ ! -d ${work_dir}/conf ]]; then
    mkdir ${work_dir}/conf
fi
ln -s ${work_dir}/configs/${project}_${env}.yaml ${work_dir}/conf/${project}.yml

#step 5 拷贝日志管理配置
cp $tmp_dir/src/store_server/logger/${project}_rotate.conf /etc/logrotate.d/${project}_rotate.conf

#step 6 install and start app
cd ${deploy_dir} && find .  -name "*.amd64.rpm" -mmin -3 -exec rpm -ivh --force {} \;
setsid service ${project} restart
