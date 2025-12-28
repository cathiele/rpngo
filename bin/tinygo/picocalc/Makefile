FS ?= littlefs
BLOCKDEV ?= machineflash
TARGET ?= pico
DATE = $(shell date +"%Y%m%d")
RELEASE = rpngo_picocalc_$(TARGET)_$(DATE).uf2

build:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -scheduler=tasks -serial=uart -o $(RELEASE)

flash:
	tinygo flash -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -serial=uart -scheduler=tasks
	

allocs:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -serial=uart -print-allocs='.*' 2>&1
