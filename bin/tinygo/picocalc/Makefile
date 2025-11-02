build:
	tinygo build -target=pico2 -scheduler=tasks -o a.uf2

flash:
	tinygo flash -target=pico2 -scheduler=tasks
	

allocs:
	tinygo build -target=pico2 -print-allocs='.*' 2>&1
