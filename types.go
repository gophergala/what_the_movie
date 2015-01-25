package main

import (
	cv `github.com/hybridgroup/go-opencv/opencv`
)

type Frame struct {
	Image   *cv.IplImage `bson:"-"`
	PosFrame int         `bson:"nframe"`
	PosMs    int         `bson:"ms"`
	Movie    Movie       `bson:"movie"`
}

type Movie struct {
	Name string          `bson:"name"`
}
