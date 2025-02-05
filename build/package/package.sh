#!/bin/sh

source ./environment.sh

# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

# bin_version=${GitTag}
bin_version=0.0.1

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rfv ../release/packages
mkdir -p ../release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ../release

for os in $os_all; do
    for arch in $arch_all; do
        bin_dir_name="${ProjectName}_${bin_version}_${os}_${arch}"
        bin_path="./packages/${ProjectName}_${bin_version}_${os}_${arch}"
        
        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./${ProjectName}_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${bin_path}
            mv ./${ProjectName}_${os}_${arch}.exe ${bin_path}/${ProjectName}.exe
        else
            if [ ! -f "./${ProjectName}_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${bin_path}
            mv ./${ProjectName}_${os}_${arch} ${bin_path}/${ProjectName}
        fi
        cp ../../LICENSE ${bin_path}
        cp -rfv ../../conf/* ${bin_path}
        
        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${bin_dir_name}.zip ${bin_dir_name}
        else
            tar -zcf ${bin_dir_name}.tar.gz ${bin_dir_name}
        fi
        cd ..
        rm -rf ${bin_path}
    done
done

cd -
