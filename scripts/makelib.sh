
name="$1"
path="$2"
import="$3"

if [[ "${path}" = "" ]];
then
    path="${name}"
    name=$(basename "${path}")
fi

if [[ "${import}" = "" ]];
then
    import="${path}"
fi

mkdir -p userlibs/$name
cd userlibs/$name
echo "module zgg_userlib_$name" > go.mod
echo "" >> go.mod
sed '1d' ../../go.mod >> go.mod
go get -u $path

srcDir=$(go mod download -json $path | grep '"Dir"' | awk -F '"' '{print $4}')
echo "srcDir -> $srcDir"
echo ../../bin/devtools export -dir "${srcDir}" -gotemplate ../../scripts/lib.tpl $import
../../bin/devtools export -dir "${srcDir}" -gotemplate ../../scripts/lib.tpl $import > $name.go
go mod tidy

echo -e "# Generated by makelib.sh. DO NOT modify

.PHONY: lib

lib:
\tgo build -buildmode=plugin -o $name.so

" > makefile
make
