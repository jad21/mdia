package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Uso:", os.Args[0], "ARCHIVO_SALIDA DIRECTORIO")
		os.Exit(1)
	}

	archivoSalida := os.Args[1]
	directorio := os.Args[2]

	// Obtener la ruta absoluta del archivo de salida
	outputFile, err := filepath.Abs(archivoSalida)
	if err != nil {
		log.Fatal(err)
	}

	// Vaciar el archivo de salida si ya existe
	err = ioutil.WriteFile(outputFile, []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Navegar al directorio especificado
	err = os.Chdir(directorio)
	if err != nil {
		log.Fatal(err)
	}

	// Encontrar todos los archivos ignorando las extensiones y carpetas especificadas
	files, err := filepath.Glob("**/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// Ignorar archivos y carpetas especificadas
		if strings.HasSuffix(file, ".pyc") ||
			strings.HasSuffix(file, ".png") ||
			strings.HasSuffix(file, ".jpg") ||
			strings.HasSuffix(file, ".jpeg") ||
			strings.HasSuffix(file, ".gif") ||
			strings.HasSuffix(file, ".svg") ||
			strings.HasSuffix(file, ".pdf") ||
			strings.HasSuffix(file, ".ico") ||
			file == "pnpm-lock.yaml" ||
			strings.HasPrefix(file, "node_modules/") ||
			strings.HasPrefix(file, "venv/") ||
			strings.HasPrefix(file, "_venv/") ||
			strings.HasPrefix(file, ".git/") ||
			strings.HasPrefix(file, "dist/") ||
			strings.HasPrefix(file, "imagenes/") ||
			strings.HasPrefix(file, "npm-locks/") {
			continue
		}

		// Obtener la ruta relativa del archivo
		relPath, err := filepath.Rel(".", file)
		if err != nil {
			log.Fatal(err)
		}

		// Obtener la extensión del archivo
		ext := filepath.Ext(file)

		// Leer el contenido del archivo
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		// Agregar al archivo de salida con el formato especificado
		fmt.Fprintf(os.Stdout, "# %s\n", relPath)
		fmt.Fprintf(os.Stdout, "```%s\n", ext[1:])
		fmt.Fprintf(os.Stdout, "%s\n", string(data))
		fmt.Fprintf(os.Stdout, "```\n\n")

		// Append al archivo de salida
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		fmt.Fprintf(f, "# %s\n", relPath)
		fmt.Fprintf(f, "```%s\n", ext[1:])
		fmt.Fprintf(f, "%s\n", string(data))
		fmt.Fprintf(f, "```\n\n")
	}
}
