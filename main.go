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
	"github.com/bmatcuk/doublestar/v4"
	ign "github.com/sabhiram/go-gitignore"
)

// multiFlag permite repetir --ignore varias veces
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func hasMeta(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

func toSlash(p string) string { return filepath.ToSlash(p) }

func main() {
	// Flags de filtro
	searchSimple := flag.String("search", "", "Subcadena a buscar en el contenido de los archivos")
	searchRegex := flag.String("search-regex", "", "Expresión regular a aplicar sobre el contenido de los archivos")
	nameSimple := flag.String("name", "", "Subcadena a buscar en la ruta/nombre de los archivos")
	nameRegex := flag.String("name-regex", "", "Expresión regular a aplicar sobre la ruta/nombre de los archivos")

	// Flag de salida
	var outFile string
	flag.StringVar(&outFile, "output", "", "Archivo de salida")
	flag.StringVar(&outFile, "o", "", "Archivo de salida (alias)")

	// --ignore (rutas o globs)
	var ignores multiFlag
	flag.Var(&ignores, "ignore", "Rutas o patrones a ignorar (repetible). Ej: --ignore src/static --ignore \"**/*.png\"")

	flag.Parse()
	dirs := flag.Args()
	if len(dirs) < 1 {
		fmt.Fprintf(os.Stderr, "Uso: %s [flags] DIRECTORIOS...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
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

	var buffer bytes.Buffer

	// Recorrer todos los directorios indicados
	for _, dir := range dirs {
		dirAbs, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("No se pudo resolver directorio %q: %v", dir, err)
		}

		// Cargar .gitignore si existe
		ignoreFile := filepath.Join(dir, ".gitignore")
		var ignorer *ign.GitIgnore
		if _, err := os.Stat(ignoreFile); err == nil {
			ignorer, _ = ign.CompileIgnoreFile(ignoreFile)
		} else {
			ignorer = ign.CompileIgnoreLines()
		}

		// Separar ignores del usuario en dos grupos:
		// 1) rutas literales (prefijo) resueltas a absolutas respecto a 'dir'
		// 2) patrones glob (doublestar) guardados tal cual y comparados contra rel y abs
		var (
			literalAbs []string
			globPats   []string
		)
		for _, ig := range ignores {
			if hasMeta(ig) {
				globPats = append(globPats, toSlash(ig))
				continue
			}
			// Ruta literal: si es relativa, interpretarla relativa a 'dir'
			p := ig
			if !filepath.IsAbs(p) {
				p = filepath.Join(dirAbs, p)
			}
			literalAbs = append(literalAbs, filepath.Clean(p))
		}

		// Helpers de ignore
		isUnderLiteral := func(absPath string) bool {
			absPath = filepath.Clean(absPath)
			for _, base := range literalAbs {
				if absPath == base || strings.HasPrefix(absPath, base+string(os.PathSeparator)) {
					return true
				}
			}
			return false
		}
		matchesGlob := func(absPath, relPath string) bool {
			a := toSlash(filepath.Clean(absPath))
			r := toSlash(filepath.Clean(relPath))
			for _, pat := range globPats {
				// Intentar contra ruta relativa al directorio raíz
				if ok, _ := doublestar.PathMatch(pat, r); ok {
					return true
				}
				// Y también contra la absoluta por si el patrón lo es
				if ok, _ := doublestar.PathMatch(pat, a); ok {
					return true
				}
			}
			return false
		}

		err = filepath.Walk(dir, func(path string, info fs.FileInfo, ferr error) error {
			if ferr != nil {
				return ferr
			}
			pathAbs, _ := filepath.Abs(path)
			relPath, _ := filepath.Rel(dirAbs, pathAbs)

			// Ignorar por --ignore (rutas o globs) o .gitignore
			if info.IsDir() {
				if isUnderLiteral(pathAbs) || matchesGlob(pathAbs, relPath) {
					return filepath.SkipDir
				}
				if ignorer != nil && ignorer.MatchesPath(path) {
					return filepath.SkipDir
				}
				return nil
			}
			if isUnderLiteral(pathAbs) || matchesGlob(pathAbs, relPath) {
				return nil
			}
			if ignorer != nil && ignorer.MatchesPath(path) {
				return nil
			}

			// Ignorar por extensiones/rutas fijas
			if strings.HasSuffix(path, ".pyc") ||
				strings.HasSuffix(path, ".png") ||
				strings.HasSuffix(path, ".jpg") ||
				strings.HasSuffix(path, ".jpeg") ||
				strings.HasSuffix(path, ".gif") ||
				strings.HasSuffix(path, ".svg") ||
				strings.HasSuffix(path, ".pdf") ||
				strings.HasSuffix(path, ".webp") ||
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

			// Incluir en buffer
			ext := filepath.Ext(path)
			relToDir, _ := filepath.Rel(dir, path)
			buffer.WriteString(fmt.Sprintf("-- `%s/%s`\n", dir, relToDir))
			buffer.WriteString(fmt.Sprintf("```%s\n", strings.TrimPrefix(ext, ".")))
			buffer.WriteString(cont + "\n```\n")
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Guardar en archivo si se indicó la flag
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
