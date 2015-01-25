package main

import (
	`fmt`
	cv `github.com/hybridgroup/go-opencv/opencv`
	`time`
)

const (
	FRAME_SKIP = 5
)

func processMovies(movies chan MovieProcessJob, framesCh chan Frame) {
	for movieJob := range movies {
		cap := cv.NewFileCapture(movieJob.Path)
		if cap == nil {
			fmt.Println(`Error: can not open video`)
			continue
		}

		start  := time.Now()
		frames := int(cap.GetProperty(cv.CV_CAP_PROP_FRAME_COUNT))
		movie  := Movie {
			Name: movieJob.Name,
		}
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
		cap.Release()
	}
}
