package provider

type cachingProvider struct {
	timesCalled int
	cache       map[string]*RoverURLsForDate
	parent      Provider
}

func NewCachedProvider(provider Provider) *cachingProvider {
	return &cachingProvider{
		cache:  map[string]*RoverURLsForDate{},
		parent: provider,
	}
}

func (c *cachingProvider) GetImagesForRover(req Request) (*RoverURLsForDate, error) {
	startKey := req.StartDate.Format("2006-01-02")

	if cacheResult, ok := c.cache[startKey]; ok {
		return cacheResult, nil
	}

	freshResult, err := c.parent.GetImagesForRover(req)
	if err != nil {
		return nil, err
	}

	c.cache[startKey] = freshResult

	return freshResult, nil
}
