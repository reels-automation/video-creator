# TODO LIST

## Video

## Subtitulos

- [x] Agregar la estructura de los subtitulos
- [x] Hacer que los subtitulos aparezcan con ffmpeg
        /*"ffmpeg -i output.mp4 -vf subtitles=subtitles.srt mysubtitledmovie.mp4"*/	


## Audio

- [] Crear la estructura del Audio
    [] Agregar audio con ffmpeg

## Imagen

- [x] Crear la estructura Imágen
- [x] Modificar su posición
- [ ] Hacerle un Resize a las imágenes


## Fixes

- [] Hacer que el video termine cuando termina la duración del audio. 
- [] 

## BuilderPattern 

Una estructura que permita crear un video completo en base a parametros.

Pensar como pasar por pipeline los distintos videos que se van creando. [NO!]

    Hay que hacer todo en una sola operación para que sea más eficiente. Buscar la forma de hacerlo en una sola operación
