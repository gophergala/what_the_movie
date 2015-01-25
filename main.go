package main

import (
	`fmt`
	`os`

	cv `github.com/hybridgroup/go-opencv/opencv`
	`time`
)

const N_THREADS = 10

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run %s videoname\n", os.Args[0])
		os.Exit(0)
	}

	filename := os.Args[1]

	cap := cv.NewFileCapture(filename)
	if cap == nil {
		panic("can not open video")
	}
	defer cap.Release()

	exitCh := make(chan struct{})
	framesCh := make(chan Frame, N_THREADS)

	for i := 0; i < N_THREADS; i++ {
		go func() {
			processFrames(framesCh)
			exitCh <- struct{}{}
		}()
	}

	now := time.Now()
	movie := Movie{Name: `test`}
	frames := int(cap.GetProperty(cv.CV_CAP_PROP_FRAME_COUNT))
	for i := 0; i < frames; i++ {
		now := time.Now()
		img := cap.QueryFrame()
		if img == nil {
			break
		}

		framesCh <- Frame {
			Image: img.Clone(),
			PosFrame: i,
			PosMs: int(cap.GetProperty(cv.CV_CAP_PROP_POS_MSEC)),
			Movie: movie,
		}
		fmt.Printf("frame: %d in %d ms.\n", i, time.Now().Sub(now)/time.Millisecond)

		if 50 == i {
			break
		}
	}
	close(framesCh)

	for i := 0; i < N_THREADS; i++ {
		<-exitCh
	}
	fmt.Printf("all frames: %d ms.\n", time.Now().Sub(now)/time.Millisecond)
}
