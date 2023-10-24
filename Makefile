# Nombre de tu programa y nombre del archivo de salida
PROGNAME = rss.exe
OUTPUT = $(PROGNAME)

# Directorio de los archivos fuente
SRC_DIR = ./cmd

# Comandos para compilar y ejecutar el programa
build:
	go build -o $(OUTPUT) $(SRC_DIR)/main.go

run: build
	./$(OUTPUT)

clean:
	rm -f ./$(OUTPUT)
