# make templates

FILE = "lib/template_bindata.go"

define Bindata_Template_HEAD
package bindata

var bindataTemplate = `
endef
export Bindata_Template_HEAD

test:
	go test -v --cover ./lib

template t:
	printf "$$Bindata_Template_HEAD" > $(FILE)
	sed "s/%/%%/g;s/package bindata/package %s/" \
		lib/bindata.go >> $(FILE)
	printf "\`" >> $(FILE)

build b:
	go build .

# generate
g:
	./bindata --resource ./resource

# example
e: g
	go build -o /tmp/example github.com/wrfly/bindata/example
	/tmp/example

.DEFAULT_GOAL := all
all: test t b e
