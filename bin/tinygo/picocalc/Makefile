FS ?= littlefs
TARGET ?= pico2

build:
	tinygo build -tags=$(FS) -target=$(TARGET) -scheduler=tasks -serial=uart -o a.uf2

flash:
	tinygo flash -tags=$(FS) -target=$(TARGET) -serial=uart -scheduler=tasks
	

allocs:
	tinygo build -tags=$(FS) -target=$(TARGET) -serial=uart -print-allocs='.*' 2>&1
