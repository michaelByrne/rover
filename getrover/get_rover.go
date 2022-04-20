package getrover

import (
	"rover/provider"
	"time"
)

type Rover struct {
	ImageURLs map[string][]string
}

type Provider interface {
	GetImagesForRover(req provider.Request) (*provider.RoverURLsForDate, error)
}

type roverGetter struct {
	provider Provider
}

func New(provider Provider) *roverGetter {
	return &roverGetter{
		provider: provider,
	}
}

func generateDates(start time.Time, numDays int) []*time.Time {
	var datesOut []*time.Time

	end := start.AddDate(0, 0, -numDays)
	for d := start; !d.Before(end.AddDate(0, 0, 1)); d = d.AddDate(0, 0, -1) {
		newDate := d
		datesOut = append(datesOut, &newDate)
	}

	return datesOut
}

func (r *roverGetter) GetRover(name string, numDays int) (*Rover, error) {
	dates := generateDates(time.Now(), numDays)
	dateMap := make(map[string][]string)
	imagesOut := &Rover{
		ImageURLs: dateMap,
	}

	for _, date := range dates {
		images, err := r.provider.GetImagesForRover(provider.Request{
			RoverName: name,
			MaxPerDay: 3,
			StartDate: date,
		})
		if err != nil {
			return nil, err
		}

		if images != nil {
			imagesOut.ImageURLs[images.Date] = images.URLs
		} else {
			imagesOut.ImageURLs[date.Format("2006-01-02")] = []string{}
		}
	}

	return imagesOut, nil
}
