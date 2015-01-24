package main

import (
	`fmt`
	`os`

	cv `github.com/hybridgroup/go-opencv/opencv`
	`time`
)

const N_THREADS = 10

func main() {
	filename := "../data/???.avi"
	if len(os.Args) == 2 {
		filename = os.Args[1]
	} else {
		fmt.Printf("Usage: go run player.go videoname\n")
		os.Exit(0)
	}

	cap := cv.NewFileCapture(filename)
	if cap == nil {
		panic("can not open video")
	}
	defer cap.Release()

	exitCh  := make(chan struct{})
	frameCh := make(chan *cv.IplImage, N_THREADS)

	for i:=0; i<N_THREADS; i++ {
		go func() {
			for img := range frameCh {
				cv.CvtColor(img, img, 40)
				img.Release()
			}
			exitCh <- struct{}{}
		}()
	}

	now := time.Now()
	frames := int(cap.GetProperty(cv.CV_CAP_PROP_FRAME_COUNT))
	for i:=0; i<frames; i++{
		now := time.Now()
		img := cap.QueryFrame()
		if img == nil {
			break
		}
		frameCh <- img.Clone()
		fmt.Printf("frame: %d in %d ms.\n", i, time.Now().Sub(now) / time.Millisecond)

		if 100 == i {
			break
		}
	}
	close(frameCh)

	for i:=0; i<N_THREADS; i++ {
		<- exitCh
	}	
	fmt.Printf("all frames: %d ms.\n", time.Now().Sub(now) / time.Millisecond)
}
