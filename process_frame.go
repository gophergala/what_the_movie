package main

/*
#cgo linux  pkg-config: opencv
#cgo darwin pkg-config: opencv
#cgo windows LDFLAGS: -lopencv_core242.dll -lopencv_imgproc242.dll -lopencv_photo242.dll -lopencv_highgui242.dll -lstdc++
#include <opencv/cv.h>
*/
import "C"

import (
	cv `github.com/hybridgroup/go-opencv/opencv`
	`unsafe`
	`math`
)

const BUCKETS = 32

var (
	rangeH   = []float32{0, 180}
	rangeS   = []float32{0, 255}
	rangeRGB = []float32{0, 255}
)

func copyHistogram(hist *C.CvHistogram) Histogram {
	values := (*cv.Mat)(hist.bins)

	ret := Histogram {
		Bins: make([]float64, values.Rows()),
	}

	for j:=0; j<len(ret.Bins); j++ {
		ret.Bins[j] = values.Get(j, 0)
	}

	return ret
}

func processFrames(frames chan Frame) {
	var img32, img1c *cv.IplImage
	var histH, histS *C.CvHistogram

	for frame := range frames {
		img := frame.Image
		if nil == img32 {
			img32  = cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 3)
			img1c  = cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 1)
			buckets := BUCKETS
			histH  = C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeH)), 1)
			histS  = C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeS)), 1)
		}

		sum := C.cvSum(unsafe.Pointer(img)).val
		max := math.Max(float64(sum[0]), float64(sum[1]))
		max  = math.Max(max, float64(sum[2]))

		r, g, b := .0, .0, .0
		if max < 0.1 {
			r, g, b = 1./3., 1./3., 1./3.
		} else {
			r = float64(sum[2])/max
			g = float64(sum[1])/max
			b = float64(sum[0])/max
		}

		C.cvConvertScale(unsafe.Pointer(img), unsafe.Pointer(img32), 1, 0)
		cv.CvtColor(img32, img32, C.CV_BGR2HSV)

		calculateHistogram(img32, img1c, histH, 1)
		calculateHistogram(img32, img1c, histS, 2)

		img.Release()

		frame.Rgb.R = r
		frame.Rgb.G = g
		frame.Rgb.B = b

		frame.Hists.H = copyHistogram(histH)
		frame.Hists.S = copyHistogram(histS)
	}

	if nil != img32 {
		img32.Release()
		img1c.Release()
		C.cvReleaseHist(&histH)
		C.cvReleaseHist(&histS)
	}
}
