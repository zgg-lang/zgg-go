
pkgs_file="$1"
if [ "${pkgs_file}" = "" ];
then
	pkgs_file='scripts/gostd_pkgs.txt'
fi

LIBDIR=stdgolibs
GO_ROOT=/usr/local/go

export CGO_ENABLED=0

envs=(
	'linux amd64'
	'linux arm64'
	'darwin amd64'
	'darwin arm64'
	'windows amd64'
)

for i in "${envs[@]}"
do
	item=($i)
	export GOOS="${item[0]}"
	export GOARCH="${item[1]}"
	while read pkg mapping
	do
		dstname="gen_$(echo ${pkg} | sed 's/\//__/g')_${GOOS}_${GOARCH}.go"
		echo "Generating ${pkg} into ${LIBDIR}/${dstname}......"
		bin/devtools export \
			--root ${GO_ROOT} \
			--go ${GO_ROOT}/bin/go \
			--gotemplate ${LIBDIR}/pkg.tpl \
			--spectypes "${mapping}" \
			${pkg} \
			> stdgolibs/${dstname}
		echo "    ${dstname} done"
	done < ${pkgs_file}
done
