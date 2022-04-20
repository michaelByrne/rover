package getrover

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rover/provider"
	"testing"
)

type fakeProvider struct {
	timesCalled int
	response *provider.RoverURLsForDate
}

func (f *fakeProvider) GetImagesForRover(req provider.Request) (*provider.RoverURLsForDate, error) {
	f.timesCalled++

	return f.response, nil
}

func TestRoverGetter_GetRover(t *testing.T) {
	t.Run("successfully gets images", func(t *testing.T) {
		response := provider.RoverURLsForDate{
			URLs: []string{"somecoolurl.jpg"},
			Date: "2015-03-15",
		}

		expectedImages := Rover{ImageURLs:map[string][]string{"2015-03-15":[]string{"somecoolurl.jpg"}}}

		fakeProvider := &fakeProvider{
			response: &response,
		}

		roverGetter := New(fakeProvider)

		actualImages, err := roverGetter.GetRover("curiosity", 3)

		require.NoError(t, err)
		assert.Equal(t, 3, fakeProvider.timesCalled)
		assert.Equal(t, &expectedImages, actualImages)
	})
}
