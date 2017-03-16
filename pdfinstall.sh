#!/bin/bash
BINDIR=/usr/local/bin
PDFPROG=${BINDIR}/wkhtmltopdf

if [ -f ${PDFPROG} ]; then
	echo "${PDFPROG} is installed."
	exit 0
fi

OSNAME=$(uname)
echo "Install for:  ${OSNAME}"

case "${OSNAME}" in
	"Darwin")
		getfile.sh ext-tools/utils/wkhtmltox-0.12.4_osx-cocoa-x86-64.pkg
		sudo installer -pkg ./wkhtmltox-0.12.4_osx-cocoa-x86-64.pkg -target /
		rm -f wkhtmltox-0.12.4_osx-cocoa-x86-64.pkg
		;;
	"Linux")
		getfile.sh ext-tools/utils/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
		tar xvf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
		;;
	*) 	echo "Unsupported operating system:  ${OSNAME}"
		exit 1
		;;
esac

