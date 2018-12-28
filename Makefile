# make templates

FILE = "lib/template_bindata.go"

define Bindata_Template_HEAD
package bindata

var bindataTemplate = `
endef
export Bindata_Template_HEAD

template t:
	printf "$$Bindata_Template_HEAD" > $(FILE)
	sed "s/%/%%/g;s/package bindata/package %s/" \
		lib/bindata.go >> $(FILE)
	printf "\`" >> $(FILE)

build b:
	go build .

.DEFAULT_GOAL := all
all: t b
