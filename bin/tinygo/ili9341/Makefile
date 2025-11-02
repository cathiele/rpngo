build:
	tinygo build -target=pico -scheduler=tasks -o a.uf2

flash:
	tinygo flash -target=pico -scheduler=tasks -monitor
	

allocs:
	tinygo build -target=pico -print-allocs='.*' 2>&1
