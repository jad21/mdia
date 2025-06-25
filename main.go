package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	argsLen := len(os.Args)
	// Acepta uno o dos parámetros: [ARCHIVO_SALIDA] DIRECTORIO
	if argsLen != 2 && argsLen != 3 {
		fmt.Fprintf(os.Stderr, "Uso: %s [ARCHIVO_SALIDA] DIRECTORIO\n", os.Args[0])
		os.Exit(1)
	}

	var archivoSalida string
	var directorio string
	if argsLen == 3 {
		archivoSalida = os.Args[1]
		directorio = os.Args[2]
	} else {
		// Sólo se proporciona el directorio; sólo guarda en portapapeles
		directorio = os.Args[1]
	}

	// Buffer para acumular la salida
	var buffer bytes.Buffer

	// Recorrer el directorio sin imprimir nada en stdout
	err := filepath.Walk(directorio, func(path string, info fs.FileInfo, err error) error {
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

	// Si se especificó archivoSalida, guardar en disco
	if archivoSalida != "" {
		outputFile, err := filepath.Abs(archivoSalida)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(outputFile, buffer.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	}

	// Copiar el contenido al portapapeles (siempre)
	if err := clipboard.WriteAll(buffer.String()); err != nil {
		log.Fatal(err)
	}
}
