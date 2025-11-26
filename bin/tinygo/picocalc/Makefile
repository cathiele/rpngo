FS ?= littlefs
BLOCKDEV ?= machineflash
TARGET ?= pico

build:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -scheduler=tasks -serial=uart

flash:
	tinygo flash -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -serial=uart -scheduler=tasks
	

allocs:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -serial=uart -print-allocs='.*' 2>&1
