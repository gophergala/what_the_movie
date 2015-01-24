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
)

var (
	buckets = 32
	rangeH  = []float32{0, 180}
	rangeS  = []float32{0, 255}
)

func processFrames(frames chan *cv.IplImage) {
	
	var img32, img1c *cv.IplImage
	var histH, histS *C.CvHistogram

	for img := range frames {
		if nil == img32 {
			img32 = cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 3)
			img1c = cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 1)
			histH = C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeH)), 1);
			histS = C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeS)), 1);
		}

		C.cvConvertScale(unsafe.Pointer(img), unsafe.Pointer(img32), 1, 0)
		cv.CvtColor(img32, img32, C.CV_BGR2HSV)

		img32.SetCOI(1)
		cv.Copy(img32, img1c, nil)
		img32.ResetROI()
		C.cvCalcHist((**C.IplImage)(unsafe.Pointer(&img1c)), histH, 0, nil)

		img32.SetCOI(2)
		cv.Copy(img32, img1c, nil)
		img32.ResetROI()
		C.cvCalcHist((**C.IplImage)(unsafe.Pointer(&img1c)), histS, 0, nil)

		img.Release()
	}

	if nil != img32 {
		img32.Release()
		img1c.Release()
	}
}
