FS ?= littlefs
BLOCKDEV ?= machineflash
TARGET ?= pico

build:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -scheduler=tasks -o a.uf2

flash:
	tinygo flash -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -scheduler=tasks
	

allocs:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -print-allocs='.*' 2>&1
