package cache

import (
	"errors"
	"hypeman/metadata"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

//MegaCache holds data caches for the various time durations...
type TimeCache struct {
	Client   *mongo.Client
	Unsorted map[string]map[string]*[]*metadata.Metadata
	Sorted   map[string]map[string]*[]*metadata.Metadata

	Incoming map[string]map[string]chan *metadata.Metadata

	//Directory will contain the location each video is allocated to, based on date
	Directory map[string]*metadata.Metadata

	Changes chan *metadata.DataChange
}

//Runs in background to unload the Changes channel
func (tc *TimeCache) UpdateCacheRoutine() {
	for {
		for len(tc.Changes) > 0 {
			change := <-tc.Changes
			md, err := tc.GetMetadata(change.Videoname)

			if err == nil {
				switch change.Category {
				case "LAUGHS":
					md.Laughs = md.Laughs + change.Delta
				case "VIEWS":
					md.Views = md.Views + change.Delta
				case "LIKES":
					md.Likes = md.Likes + change.Delta
				case "DISLIKES":
					md.Dislikes = md.Dislikes + change.Delta
				}
			}
		}

		time.Sleep(time.Second * time.Duration(30))
	}
}

//if cat == """, update every cache. otherwise, just update the category specified
func (tc *TimeCache) AddToCache(md *metadata.Metadata, cat string) {
	time := md.Bucket

	tc.Directory[md.Videoname] = md

	if cat == "" {
		for _, category := range categories {
			tc.Incoming[time][category] <- md
		}
	} else {
		tc.Incoming[time][cat] <- md
	}
}

func (tc *TimeCache) GetVideos(time string, category string, needed int, startIndex int) (*[]*metadata.Metadata, error) {
	if _, ok := tc.Sorted[time]; !ok {
		return nil, errors.New("Invalid [time]")
	}

	if _, ok := tc.Sorted[time][category]; !ok {
		return nil, errors.New("Invalid [category]")
	}

	sorted := *tc.Sorted[time][category]

	maxIndex := len(sorted) - 1
	if maxIndex+1 < needed {
		return &sorted, nil
	}

	//bound startIndex to be >= 0
	if startIndex < 0 {
		startIndex = 0
	}

	//ensure startindex is proper
	if startIndex > maxIndex {
		return nil, errors.New("Start index exceeds slice length!")
	}

	arr := []*metadata.Metadata{}

	for i := startIndex; i < startIndex+needed && i <= maxIndex; i++ {
		arr = append(arr, sorted[i])
	}

	return &arr, nil
}

var categories = []string{"NEWEST", "FUNNIEST", "MOSTLIKES", "MOSTDISLIKES", "MOSTSEEN"}

func (tc *TimeCache) GetMetadata(video string) (*metadata.Metadata, error) {
	if v, ok := tc.Directory[video]; ok {
		return v, nil
	}

	data, err := tc.GetVideoFromDB(video)

	if err == nil {
		tc.AddToCache(data, "")

		return data, err
	}

	return nil, errors.New("Metadata not in cache or database!")
}

func NewTimeCache(client *mongo.Client) *TimeCache {
	tc := &TimeCache{
		Client:    client,
		Directory: make(map[string]*metadata.Metadata),
		Unsorted:  make(map[string]map[string]*[]*metadata.Metadata),
		Sorted:    make(map[string]map[string]*[]*metadata.Metadata),
		Incoming:  make(map[string]map[string]chan *metadata.Metadata),
		Changes:   make(chan *metadata.DataChange, 1000),
	}

	//Dynamically create & populate all the buckets that will populate this cache
	times := []string{"TODAY", "THISWEEK", "THISMONTH", "THISYEAR", "ALLTIME"}

	for _, t := range times {
		tc.Unsorted[t] = make(map[string]*[]*metadata.Metadata)
		tc.Sorted[t] = make(map[string]*[]*metadata.Metadata)
		tc.Incoming[t] = make(map[string]chan *metadata.Metadata, 500)

		x, _ := tc.GetBucketVideosFromDB(t)

		if x == nil {
			x = &[]*metadata.Metadata{}
		}

		for _, video := range *x {
			tc.Directory[video.Videoname] = video
		}

		for _, c := range categories {
			tc.Unsorted[t][c] = x
			tc.Incoming[t][c] = make(chan *metadata.Metadata, 500) //can handle 500 incoming videos every 6 minutes
		}
	}
	tc.ReSort()

	return tc
}

var lastWeek = time.Now().AddDate(0, 0, -7).Unix()
var lastMonth = time.Now().AddDate(0, -1, 0).Unix()
var lastYear = time.Now().AddDate(-1, 0, 0).Unix()
var yesterday = time.Now().AddDate(0, 0, -1).Unix()

//Input is a video metadata and the output is what category it belongs in!
func TimeValidate(md *metadata.Metadata) string {
	ti := int64(md.Date)

	if ti > yesterday {
		return "TODAY"
	}

	if ti > lastWeek {
		return "THISWEEK"
	}

	if ti > lastMonth {
		return "THISMONTH"
	}

	if ti > lastYear {
		return "THISYEAR"
	}

	return "ALLTIME"
}

func updateTimes() {
	lastWeek = time.Now().AddDate(0, 0, -7).Unix()
	lastMonth = time.Now().AddDate(0, -1, 0).Unix()
	lastYear = time.Now().AddDate(-1, 0, 0).Unix()
	yesterday = time.Now().AddDate(0, 0, -1).Unix()
}

func (tc *TimeCache) PerpetualSortRoutine(min int) {
	for {
		tc.ReSort()
		time.Sleep(time.Minute * time.Duration(min))
	}
}

func (tc *TimeCache) ReSort() {
	times := []string{"TODAY", "THISWEEK", "THISMONTH", "THISYEAR", "ALLTIME"}

	updateTimes()

	for _, time := range times {

		iteration := 0
		for _, category := range categories {

			mds := []*metadata.Metadata{}

			if n, ok := tc.Unsorted[time][category]; !ok || n == nil {
				tc.Unsorted[time][category] = &[]*metadata.Metadata{}
			}

			for _, md := range *tc.Unsorted[time][category] {
				trueTime := TimeValidate(md)

				if trueTime == time {
					mds = append(mds, md)
				} else { //if time needs alterations
					md.Bucket = trueTime

					if iteration == 0 {
						tc.UpdateDBTime(md)
					}
					tc.AddToCache(md, category)
				}
			}

			iteration += 1
			//add from incoming channels
			for len(tc.Incoming[time][category]) > 0 {
				elem := <-tc.Incoming[time][category]
				mds = append(mds, elem)
			}

			switch category {
			case "NEWEST":
				sort.SliceStable(mds, func(i, j int) bool {
					return mds[i].Date > mds[i].Date
				})
			case "FUNNIEST":
				sort.SliceStable(mds, func(i, j int) bool {
					return mds[i].Laughs < mds[i].Laughs
				})
			case "MOSTLIKES":
				sort.SliceStable(mds, func(i, j int) bool {
					return mds[i].Likes < mds[i].Likes
				})
			case "MOSTDISLIKES":
				sort.SliceStable(mds, func(i, j int) bool {
					return mds[i].Dislikes < mds[i].Dislikes
				})
			case "MOSTVIEWS":
				sort.SliceStable(mds, func(i, j int) bool {
					return mds[i].Views < mds[i].Views
				})
			}

			tc.Sorted[time][category] = &mds
			tc.Unsorted[time][category] = &mds
		}
	}
}
