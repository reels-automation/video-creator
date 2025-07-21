# TODO LIST

## Video

## Subtitulos

- [x] Agregar la estructura de los subtitulos
- [x] Hacer que los subtitulos aparezcan con ffmpeg
        /*"ffmpeg -i output.mp4 -vf subtitles=subtitles.srt mysubtitledmovie.mp4"*/	

## Audio

- [x] Crear la estructura del Audio
    [x] Agregar audio con ffmpeg

## Imagen

- [x] Crear la estructura Imágen
- [x] Modificar su posición
- [x] Hacerle un Resize a las imágenes
- [ ] Eliminar el hardcodeo de las imágenes

## Fixes

- [x] Hacer que el video termine cuando termina la duración del audio. 
- [] 

## BuilderPattern 

Una estructura que permita crear un video completo en base a parametros.

Pensar como pasar por pipeline los distintos videos que se van creando. [NO!]

    Hay que hacer todo en una sola operación para que sea más eficiente. Buscar la forma de hacerlo en una sola operación

## Hacer productor de Kafka

- [ ] Hacer un productor de Kafka que suba los videos al bucket y los publique en la API.

## main.go

Limpiar y separar main.go !