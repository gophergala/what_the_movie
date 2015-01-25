package main

import (
	`os`
	mgo `gopkg.in/mgo.v2`
	`gopkg.in/mgo.v2/bson`
	log `github.com/cihub/seelog`
	`time`
)

const (
	N_FRAME_THREADS = 10
	N_MOVIE_THREADS = 5
	JOBS_COLLECTION = `jobs`
	FRAMES_COLLECTION = `frames`
	MOVIES_COLLECTION = `movies`
)

func main() {
	defer log.Flush()
	if 3 != len(os.Args) {
		log.Errorf(`Usage: %s dbhost dbname`, os.Args[0])
		return
	}

	// Database connection.
	log.Infof(`Connecting to database %s at address %s.`, os.Args[2], os.Args[1])
	mgoSession, err := mgo.Dial(os.Args[1])
	if err != nil {
		log.Errorf(`Error connecting to MongoDB: %v.`, err)
		return
	}

	log.Infof(`Connected!`)
	dbConn := mgoSession.DB(os.Args[2])

	moviesCh := make(chan MovieProcessJob, N_MOVIE_THREADS)
	framesCh := make(chan Frame, N_FRAME_THREADS)
	defer close(moviesCh)
	defer close(framesCh)

	for i := 0; i < N_MOVIE_THREADS; i++ {
		go processMovies(moviesCh, framesCh, mgoSession.Copy(), os.Args[2])
	}
	for i := 0; i < N_FRAME_THREADS; i++ {
		go processFrames(framesCh, mgoSession.Copy(), os.Args[2])
	}

	for {
		log.Infof(`Querying...`)
		jobsColl := dbConn.C(JOBS_COLLECTION)
		jobsIter := jobsColl.Find(bson.M{`processed`: false}).Tail(-1)
		var movieJob MovieProcessJob
		for jobsIter.Next(&movieJob) {
			moviesCh <- movieJob
		}
		if err := jobsIter.Err(); nil != err {
			log.Errorf(`Error: %v.`, err)
		}
		jobsIter.Close()
		time.Sleep(time.Second)
	}
}
