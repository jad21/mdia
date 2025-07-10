# MDia (Markdown para IA)

Este es un script en Go que se puede ejecutar desde la línea de comandos. Permite generar un archivo de salida con el contenido de todos los archivos de un directorio, excluyendo ciertos tipos de archivos y carpetas, y copiar dicho contenido al portapapeles.

## Instalación

1. Clona el repositorio:

   ```bash
   git clone https://github.com/jad21/mdia.git
   ```

2. Instala el binario con Go:

   ```bash
   go install github.com/jad21/mdia@latest
   ```

3. Instala el binario usando Goblin:

   ```bash
   curl -sf https://goblin.run/github.com/jad21/mdia@v0.1.0 | sh
   ```

## Uso

El comportamiento varía según el número de parámetros:

* **Solo directorio** (1 parámetro): copia todo el contenido procesado al portapapeles.

  ```bash
  mdia DIRECTORIO
  ```

* **Archivo de salida y directorio** (2 parámetros): guarda el contenido en el archivo especificado y también lo copia al portapapeles.

  ```bash
  mdia ARCHIVO_SALIDA DIRECTORIO
  ```

Donde:

* `DIRECTORIO` es el directorio del cual se leerán los archivos.
* `ARCHIVO_SALIDA` es la ruta del archivo de salida que se generará (opcional si solo se desea copiar al portapapeles).

El script procesará cada archivo y formateará la salida así:

````
# ruta/relativa/del/archivo.ext
```ext
Contenido del archivo
````

## Características

* Ignora archivos con las siguientes extensiones: `.pyc`, `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, `.pdf`, `.ico`.
* Ignora los siguientes archivos y carpetas: `pnpm-lock.yaml`, `node_modules/`, `venv/`, `_venv/`, `.git/`, `dist/`, `imagenes/`, `npm-locks/`.
* Obtiene la ruta relativa de cada archivo.
* Copia el resultado al portapapeles.

## Contribución

Si encuentras algún problema o tienes sugerencias de mejora, no dudes en abrir un issue o enviar un pull request.

## Licencia

Este proyecto se distribuye bajo la [Licencia MIT](LICENSE).
