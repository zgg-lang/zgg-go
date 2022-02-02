
pkgs_file="$1"
if [ "${pkgs_file}" = "" ];
then
	pkgs_file='scripts/gostd_pkgs.txt'
fi

LIBDIR=stdgolibs
GO_ROOT=/usr/local/go

export CGO_ENABLED=0

export GOOS=linux
export GOARCH=amd64

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

export GOOS=darwin
export GOARCH=amd64

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
		> ${LIBDIR}/${dstname}
	echo "    ${dstname} done"
done < ${pkgs_file}

export GOOS=darwin
export GOARCH=arm64

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
		> ${LIBDIR}/${dstname}
	echo "    ${dstname} done"
done < ${pkgs_file}

export GOOS=windows
export GOARCH=amd64

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
		> ${LIBDIR}/${dstname}
	echo "    ${dstname} done"
done < ${pkgs_file}
