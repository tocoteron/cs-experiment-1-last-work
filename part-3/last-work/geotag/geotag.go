package geotag

import (
	"cs-experiment-1/part-3/last-work/csvutil"
	"os"
	"strconv"
)

// GeoTag is corresponding the db table 'geotag'
type GeoTag struct {
	ID        int
	Time      string
	Latitude  float64
	Longitude float64
	URL       string
}

func UnmarshalGeoTag(data []string) (GeoTag, error) {
	id, err := strconv.Atoi(data[0])
	if err != nil {
		return GeoTag{}, err
	}

	latitude, err := strconv.ParseFloat(data[2], 64)
	if err != nil {
		return GeoTag{}, err
	}

	longitude, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		return GeoTag{}, err
	}

	geotag := GeoTag{
		ID:        id,
		Time:      data[1],
		Latitude:  latitude,
		Longitude: longitude,
		URL:       data[4],
	}

	return geotag, nil
}

func UnmarshalGeoTags(data [][]string) ([]GeoTag, error) {
	var geotags []GeoTag

	for i := 0; i < len(data); i++ {
		geotag, err := UnmarshalGeoTag(data[i])
		if err != nil {
			return nil, err
		}
		geotags = append(geotags, geotag)
	}

	return geotags, nil
}

func ReadGeoTagsFromCSV(path string, capacity int, buffsize int) ([]GeoTag, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	geotags := make([]GeoTag, 0, capacity)

	for record := range csvutil.AsyncReadCSV(reader, buffsize) {
		geotag, err := UnmarshalGeoTag(record)
		if err != nil {
			return nil, err
		}

		geotags = append(geotags, geotag)
	}

	return geotags, nil
}
