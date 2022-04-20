package provider

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"path"
	"testing"
	"time"
)

type fakeHTTPClient struct {
	timesCalled int
	t           *testing.T
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	f.timesCalled++

	responseBody := ioutil.NopCloser(bytes.NewReader([]byte(loadFile(f.t, "test_response.json"))))
	return &http.Response{
		StatusCode: 200,
		Body:       responseBody,
	}, nil
}

func TestProvider_GetImagesForRover(t *testing.T) {
	t.Run("successfully gets three images per day", func(t *testing.T) {
		client := New(BaseURL, WithHTTPClient(&fakeHTTPClient{t: t}))
		day := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

		images, err := client.GetImagesForRover(Request{
			RoverName: "curiosity",
			MaxPerDay: 3,
			StartDate: &day,
		})

		expected := RoverURLsForDate{URLs: []string{"http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480140156EDR_F0450852FHAZ00323M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FLB_480154331EDR_F0450852FHAZ00190M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480154331EDR_F0450852FHAZ00190M_.JPG"}, Date: "2015-03-20"}

		require.NoError(t, err)
		assert.Equal(t, &expected, images)
	})

	t.Run("successfully gets five images per day", func(t *testing.T) {
		client := New(BaseURL, WithHTTPClient(&fakeHTTPClient{t: t}))
		day := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

		images, err := client.GetImagesForRover(Request{
			RoverName: "curiosity",
			MaxPerDay: 5,
			StartDate: &day,
		})

		expected := RoverURLsForDate{URLs: []string{"http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480140156EDR_F0450852FHAZ00323M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FLB_480154331EDR_F0450852FHAZ00190M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480154331EDR_F0450852FHAZ00190M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FLB_480140156EDR_F0450852FHAZ00323M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/rcam/RRB_480140189EDR_F0450852RHAZ00323M_.JPG"}, Date: "2015-03-20"}

		require.NoError(t, err)
		assert.Equal(t, &expected, images)
	})
}

func loadFile(t *testing.T, filename string) string {
	if filename == "" {
		return ""
	}

	fileLocation := path.Join("testdata", filename)
	fileBytes, err := ioutil.ReadFile(fileLocation)
	assert.NoErrorf(t, err, "error opening file: %q", filename)

	fileString := string(fileBytes)
	assert.NotEmptyf(t, fileString, "file is empty: %q", filename)

	return string(fileBytes)
}
