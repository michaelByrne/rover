package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const BaseURL = "https://api.nasa.gov/mars-photos/api/v1/rovers/%s/photos?earth_date=%s&api_key=%s"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Provider interface {
	GetImagesForRover(req Request) (*RoverURLsForDate, error)
}

type provider struct {
	apiKey     string
	httpClient HTTPClient
	baseURL    string
}

type Request struct {
	RoverName string
	MaxPerDay int
	StartDate *time.Time
}

type Response struct {
	Photos []struct {
		ID     int `json:"id"`
		Sol    int `json:"sol"`
		Camera struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			RoverID  int    `json:"rover_id"`
			FullName string `json:"full_name"`
		} `json:"camera"`
		ImgSrc    string `json:"img_src"`
		EarthDate string `json:"earth_date"`
		Rover     struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			LandingDate string `json:"landing_date"`
			LaunchDate  string `json:"launch_date"`
			Status      string `json:"status"`
		} `json:"rover"`
	} `json:"photos"`
}

type RoverURLsForDate struct {
	URLs []string
	Date string
}

type Option func(p *provider)

func New(baseURL string, options ...Option) *provider {
	p := &provider{
		baseURL: baseURL,
	}

	for _, opt := range options {
		opt(p)
	}

	return p
}

func WithHTTPClient(client HTTPClient) Option {
	return func(p *provider) {
		p.httpClient = client
	}
}

func WithAPIKey(key string) Option {
	return func(p *provider) {
		p.apiKey = key
	}
}

func (p *provider) GetImagesForRover(req Request) (*RoverURLsForDate, error) {
	if p.httpClient == nil {
		p.httpClient = &http.Client{}
	}

	apiKey := "DEMO_KEY"

	if p.apiKey != "" {
		apiKey = p.apiKey
	}

	url := fmt.Sprintf(p.baseURL, req.RoverName, req.StartDate.Format("2006-01-02"), apiKey)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	response, err := p.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var roverResponse Response
	err = json.Unmarshal(body, &roverResponse)
	if err != nil {
		return nil, err
	}

	return roverResponse.toRoverImages(req.MaxPerDay), nil
}

func (r *Response) toRoverImages(maxPerDay int) *RoverURLsForDate {
	if len(r.Photos) == 0 {
		return nil
	}

	var outURLs []string

	for dex, pic := range r.Photos {
		if dex > maxPerDay-1 {
			break
		}

		outURLs = append(outURLs, pic.ImgSrc)
	}

	return &RoverURLsForDate{
		URLs: outURLs,
		Date: r.Photos[0].EarthDate,
	}
}
