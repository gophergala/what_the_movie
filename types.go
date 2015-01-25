package main

import (
	cv `github.com/hybridgroup/go-opencv/opencv`
)


type Movie struct {
	Name string `bson:"name"`
}

type RGB struct {
	R float64 `bson:"r"`
	G float64 `bson:"g"`
	B float64 `bson:"b"`
}

type Histogram struct {
	Bins []float64 `bson:"bins"`
}

type Histograms struct {
	H Histogram `bson:"h"`
	S Histogram `bson:"s"`
}

type Frame struct {
	Image   *cv.IplImage `bson:"-"`
	PosFrame int         `bson:"nframe"`
	PosMs    int         `bson:"ms"`
	Movie    Movie       `bson:"movie"`
	Rgb      RGB         `bson:"rgb"`
	Hists    Histograms  `bson:"hists"`
}
