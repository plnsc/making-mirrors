#!/usr/bin/env bash
magick originals/pe-de-acerola-16x16.xcf -resize 16x16 -background none -gravity center -extent 16x16 static/favicon.ico
magick originals/pe-de-acerola-32x32.xcf -resize 32x32 -background none -gravity center -extent 32x32 static/favicon-32x32.ico
magick originals/pe-de-acerola-32x32.xcf -define webp:lossless=true -resize 32x32 -background none -gravity center -extent 32x32 static/pe-de-acerola-32x32.webp
