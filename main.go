package main

import (
	`fmt`
	`os`

	cv `github.com/hybridgroup/go-opencv/opencv`
	`time`
)

const N_THREADS = 10
const FRAME_SKIP = 5

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

	start := time.Now()
	movie := Movie{Name: `test`}
	frames := int(cap.GetProperty(cv.CV_CAP_PROP_FRAME_COUNT))
	for i := 0; i < frames; i++ {
		img := cap.QueryFrame()
		if img == nil {
			break
		}

		fmt.Printf("Processing frame %d (%.2f%%). %.2f fps.\n", i, float32(i)*100./float32(frames), (float32(i*1e9))/float32(time.Now().Sub(start).Nanoseconds()))

		framesCh <- Frame {
			Image: img.Clone(),
			PosFrame: i,
			PosMs: int(cap.GetProperty(cv.CV_CAP_PROP_POS_MSEC)),
			Movie: movie,
		}

		// Skip N-1 frames
		for j:=1; j<FRAME_SKIP; j++ {
			cap.GrabFrame()
			i++
		}
	}
	close(framesCh)

	for i := 0; i < N_THREADS; i++ {
		<-exitCh
	}
	fmt.Printf("All frames: %.2f s.\n", time.Now().Sub(start)/time.Second)
}
