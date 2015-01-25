package main

const (
	N_FRAME_THREADS = 10
	N_MOVIE_THREADS = 5
)

func main() {
	moviesCh := make(chan MovieProcessJob, N_MOVIE_THREADS)
	framesCh := make(chan Frame, N_FRAME_THREADS)
	defer close(moviesCh)
	defer close(framesCh)

	for i := 0; i < N_MOVIE_THREADS; i++ {
		go processMovies(moviesCh, framesCh)
	}
	for i := 0; i < N_FRAME_THREADS; i++ {
		go processFrames(framesCh)
	}

	<- (chan struct{})(nil)
}
