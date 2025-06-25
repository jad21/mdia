# Dockia

Este es un script en Go que se puede ejecutar desde la línea de comandos. Permite generar un archivo de salida con el contenido de todos los archivos de un directorio, excluyendo ciertos tipos de archivos y carpetas.

## Instalación

1. Clona el repositorio:

```
git clone https://github.com/jad21/dockia.git
```

2. Instala el binario:

```
go install github.com/jad21/dockia
```

## Uso

Para ejecutar el script, utiliza el siguiente comando:

```
dockia ARCHIVO_SALIDA DIRECTORIO
```

Donde:

- `ARCHIVO_SALIDA` es la ruta del archivo de salida que se generará.
- `DIRECTORIO` es el directorio del cual se leerán los archivos.

El script generará un archivo de salida con el siguiente formato:

```
# ruta/relativa/del/archivo.ext
```
```ext
Contenido del archivo
```

## Características

- Ignora archivos con las siguientes extensiones: `.pyc`, `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, `.pdf`, `.ico`.
- Ignora los siguientes archivos y carpetas: `pnpm-lock.yaml`, `node_modules/`, `venv/`, `_venv/`, `.git/`, `dist/`, `imagenes/`, `npm-locks/`.
- Obtiene la ruta relativa de cada archivo.
- Muestra el contenido de cada archivo con el formato de código correspondiente a su extensión.

## Contribución

Si encuentras algún problema o tienes sugerencias de mejora, no dudes en abrir un issue o enviar un pull request.

## Licencia

Este proyecto se distribuye bajo la [Licencia MIT](LICENSE).
