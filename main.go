package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Uso: %s ARCHIVO_SALIDA DIRECTORIO\n", os.Args[0])
		os.Exit(1)
	}

	archivoSalida := os.Args[1]
	directorio := os.Args[2]

	// Obtener la ruta absoluta del archivo de salida
	outputFile, err := filepath.Abs(archivoSalida)
	if err != nil {
		log.Fatal(err)
	}

	// Buffer para acumular la salida
	var buffer bytes.Buffer

	// Recorrer el directorio sin imprimir nada en stdout
	err = filepath.Walk(directorio, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Ignorar archivos y carpetas especificadas
		if strings.HasSuffix(path, ".pyc") ||
			strings.HasSuffix(path, ".png") ||
			strings.HasSuffix(path, ".jpg") ||
			strings.HasSuffix(path, ".jpeg") ||
			strings.HasSuffix(path, ".gif") ||
			strings.HasSuffix(path, ".svg") ||
			strings.HasSuffix(path, ".pdf") ||
			strings.HasSuffix(path, ".ico") ||
			path == "pnpm-lock.yaml" ||
			strings.HasPrefix(path, "node_modules/") ||
			strings.HasPrefix(path, "venv/") ||
			strings.HasPrefix(path, "_venv/") ||
			strings.HasPrefix(path, ".git/") ||
			strings.HasPrefix(path, "dist/") ||
			strings.HasPrefix(path, "imagenes/") ||
			strings.HasPrefix(path, "npm-locks/") {
			return nil
		}

		// Obtener la ruta relativa
		relPath, err := filepath.Rel(directorio, path)
		if err != nil {
			return err
		}

		// Leer contenido
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)

		// Escribir al buffer
		fmt.Fprintf(&buffer, "-- %s\n", relPath)
		fmt.Fprintf(&buffer, "```%s\n", strings.TrimPrefix(ext, "."))
		fmt.Fprintf(&buffer, "%s\n", data)
		fmt.Fprintf(&buffer, "```\n\n")

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Guardar el contenido acumulado en el archivo de salida
	if err := os.WriteFile(outputFile, buffer.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
