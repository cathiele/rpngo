FS ?= littlefs
TARGET ?= pico2

build:
	tinygo build -tags=$(FS) -target=$(TARGET) -scheduler=tasks -o a.uf2

flash:
	tinygo flash -tags=$(FS) -target=$(TARGET) -scheduler=tasks
	

allocs:
	tinygo build -tags=$(FS) -target=$(TARGET) -print-allocs='.*' 2>&1
