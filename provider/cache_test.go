package provider

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCacheProvider_GetImagesForRover(t *testing.T) {
	t.Run("successfully gets images for three days from cache", func(t *testing.T) {
		client := New(BaseURL, WithHTTPClient(&fakeHTTPClient{t: t}))
		cacheClient := NewCachedProvider(client)
		cacheClient.cache["2015-03-15"] = &RoverURLsForDate{
			URLs: []string{"http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480140156EDR_F0450852FHAZ00323M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FLB_480154331EDR_F0450852FHAZ00190M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480154331EDR_F0450852FHAZ00190M_.JPG"},
			Date: "2015-03-15",
		}
		day := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

		images, err := cacheClient.GetImagesForRover(Request{
			RoverName: "curiosity",
			MaxPerDay: 3,
			StartDate: &day,
		})

		expected := RoverURLsForDate{URLs: []string{"http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480140156EDR_F0450852FHAZ00323M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FLB_480154331EDR_F0450852FHAZ00190M_.JPG", "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/00931/opgs/edr/fcam/FRB_480154331EDR_F0450852FHAZ00190M_.JPG"}, Date: "2015-03-15"}

		require.NoError(t, err)
		assert.Equal(t, &expected, images)
		assert.Zero(t, cacheClient.timesCalled)
	})
}
