build:
	if [ ! -d ./bin ]; then mkdir -p ./bin; fi
	gcc ./src/machete/mach.c -o ./bin/mach_c
clean:
	rm ./bin/*
