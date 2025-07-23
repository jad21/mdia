// main.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	ign "github.com/sabhiram/go-gitignore"
)

func main() {
	// Flags de filtro
	searchSimple := flag.String("search", "", "Subcadena a buscar en el contenido de los archivos")
	searchRegex := flag.String("search-regex", "", "Expresión regular a aplicar sobre el contenido de los archivos")
	nameSimple := flag.String("name", "", "Subcadena a buscar en la ruta/nombre de los archivos")
	nameRegex := flag.String("name-regex", "", "Expresión regular a aplicar sobre la ruta/nombre de los archivos")

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 && len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Uso: %s [ARCHIVO_SALIDA] DIRECTORIO [flags]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Determinar salida y directorio
	var outFile, dir string
	if len(args) == 2 {
		outFile = args[0]
		dir = args[1]
	} else {
		dir = args[0]
	}

	// Preparar regex si se especificaron
	var (
		reContent *regexp.Regexp
		reName    *regexp.Regexp
		err       error
	)
	if *searchRegex != "" {
		if reContent, err = regexp.Compile(*searchRegex); err != nil {
			log.Fatalf("Regex de contenido inválida: %v", err)
		}
	}
	if *nameRegex != "" {
		if reName, err = regexp.Compile(*nameRegex); err != nil {
			log.Fatalf("Regex de nombres inválida: %v", err)
		}
	}

	// Cargar .gitignore
	ignoreFile := filepath.Join(dir, ".gitignore")
	var ignorer *ign.GitIgnore
	if _, err := os.Stat(ignoreFile); err == nil {
		ignorer, _ = ign.CompileIgnoreFile(ignoreFile)
	} else {
		ignorer = ign.CompileIgnoreLines()
	}

	var buffer bytes.Buffer
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, ferr error) error {
		if ferr != nil {
			return ferr
		}
		if info.IsDir() {
			if ignorer != nil && ignorer.MatchesPath(path) {
				return filepath.SkipDir
			}
			return nil
		}
		// Aplicar ignore por gitignore y extensiones fijas
		if (ignorer != nil && ignorer.MatchesPath(path)) ||
			strings.HasSuffix(path, ".pyc") ||
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

		// Filtro por nombre
		if *nameSimple != "" && !strings.Contains(path, *nameSimple) {
			return nil
		}
		if reName != nil && !reName.MatchString(path) {
			return nil
		}

		// Leer contenido
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		cont := string(data)

		// Filtro por contenido
		if *searchSimple != "" && !strings.Contains(cont, *searchSimple) {
			return nil
		}
		if reContent != nil && !reContent.MatchString(cont) {
			return nil
		}

		// Si pasa todos los filtros, lo incluimos
		ext := filepath.Ext(path)
		fmt.Fprintf(&buffer, "-- `%s`\n", strings.TrimPrefix(path, dir+"/"))
		fmt.Fprintf(&buffer, "```%s\n", strings.TrimPrefix(ext, "."))
		fmt.Fprintf(&buffer, "%s\n```\n\n", cont)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Guardar en archivo si se indicó
	if outFile != "" {
		abs, err := filepath.Abs(outFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(abs, buffer.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	}

	// Copiar al portapapeles
	if err := clipboard.WriteAll(buffer.String()); err != nil {
		log.Fatal(err)
	}
}
