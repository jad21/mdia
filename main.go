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
	ign "github.com/sabhiram/go-gitignore"
)

func main() {
	argsLen := len(os.Args)
	if argsLen != 2 && argsLen != 3 {
		fmt.Fprintf(os.Stderr, "Uso: %s [ARCHIVO_SALIDA] DIRECTORIO\n", os.Args[0])
		os.Exit(1)
	}

	var archivoSalida, directorio string
	if argsLen == 3 {
		nArchivo := os.Args[1]
		archivoSalida = nArchivo
		directorio = os.Args[2]
	} else {
		directorio = os.Args[1]
	}

	// Cargar .gitignore si existe
	ignoreFile := filepath.Join(directorio, ".gitignore")
	var ignorer *ign.GitIgnore
	if _, err := os.Stat(ignoreFile); err == nil {
		ignorer, _ = ign.CompileIgnoreFile(ignoreFile)
	} else {
		// Ningún .gitignore presente o no accesible
		ignorer = ign.CompileIgnoreLines()
	}

	var buffer bytes.Buffer

	err := filepath.Walk(directorio, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Saltar directorios
		if info.IsDir() {
			// Ignorar directorios según gitignore
			if ignorer != nil && ignorer.MatchesPath(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Ignorar según .gitignore
		if ignorer != nil && ignorer.MatchesPath(path) {
			return nil
		}

		// Ignorar extensiones o patrones fijos
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

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		fmt.Fprintf(&buffer, "-- `%s`\n", strings.TrimSpace(path))
		fmt.Fprintf(&buffer, "```%s\n", strings.TrimPrefix(ext, "."))
		fmt.Fprintf(&buffer, "%s\n```\n\n", strings.TrimSpace(string(data)))

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if archivoSalida != "" {
		outputFile, err := filepath.Abs(archivoSalida)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(outputFile, buffer.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	}

	if err := clipboard.WriteAll(buffer.String()); err != nil {
		log.Fatal(err)
	}
}
