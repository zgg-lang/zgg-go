
pkgs_file="$1"
if [ "${pkgs_file}" = "" ];
then
	pkgs_file='scripts/gostd_pkgs.txt'
fi

export CGO_ENABLED=0

export GOOS=linux
export GOARCH=amd64

while read pkg mapping
do
	dstname="gen_$(echo ${pkg} | sed 's/\//__/g')_${GOOS}_${GOARCH}.go"
	echo "Generating ${pkg} into stdgolibs/${dstname}......"
	bin/devtools export \
		--root /usr/local/go \
		--gotemplate stdgolibs/pkg.tpl \
		--spectypes "${mapping}" \
		${pkg} \
		> stdgolibs/${dstname}
	echo "    ${dstname} done"
done < ${pkgs_file}

export GOOS=darwin
export GOARCH=amd64

while read pkg mapping
do
	dstname="gen_$(echo ${pkg} | sed 's/\//__/g')_${GOOS}_${GOARCH}.go"
	echo "Generating ${pkg} into src/zgg/stdgolibs/${dstname}......"
	bin/devtools export \
		--root /usr/local/go \
		--gotemplate src/zgg/stdgolibs/pkg.tpl \
		--spectypes "${mapping}" \
		${pkg} \
		> src/zgg/stdgolibs/${dstname}
	echo "    ${dstname} done"
done < ${pkgs_file}

export GOOS=darwin
export GOARCH=arm64

while read pkg mapping
do
	dstname="gen_$(echo ${pkg} | sed 's/\//__/g')_${GOOS}_${GOARCH}.go"
	echo "Generating ${pkg} into src/zgg/stdgolibs/${dstname}......"
	bin/devtools export \
		--root /usr/local/go \
		--gotemplate src/zgg/stdgolibs/pkg.tpl \
		--spectypes "${mapping}" \
		${pkg} \
		> src/zgg/stdgolibs/${dstname}
	echo "    ${dstname} done"
done < ${pkgs_file}
