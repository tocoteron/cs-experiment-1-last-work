package geotag

import (
	"cs-experiment-1/part-3/last-work/csvutil"
	"os"
	"strconv"
	"time"
)

// GeoTag is corresponding the db table 'geotag'
type GeoTag struct {
	ID        int
	Time      int32
	Latitude  float64
	Longitude float64
	URL       string
}

// Datetime is raw datetime info
func (geotag *GeoTag) Datetime() string {
	t := time.Unix(int64(geotag.Time), 0)
	return t.Format("2006-01-02 15:04:05")
}

func unixtime(datetime string, timezone *time.Location) (int32, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, timezone)
	if err != nil {
		return 0, err
	}

	return int32(t.Unix()), nil
}

// UnmarshalGeoTag converts []string to GeoTag
func UnmarshalGeoTag(data []string, timezone *time.Location) (GeoTag, error) {
	id, err := strconv.Atoi(data[0])
	if err != nil {
		return GeoTag{}, err
	}

	datetime, err := unixtime(data[1], timezone)
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
		Time:      datetime,
		Latitude:  latitude,
		Longitude: longitude,
		URL:       data[4],
	}

	return geotag, nil
}

// ReadGeoTagsFromCSV marshal []GeoTag from CSV file corresponding the path
func ReadGeoTagsFromCSV(path string, capacity int, buffsize int) ([]GeoTag, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	geotags := make([]GeoTag, 0, capacity)

	for record := range csvutil.AsyncReadCSV(reader, buffsize) {
		geotag, err := UnmarshalGeoTag(record, jst)
		if err != nil {
			return nil, err
		}

		geotags = append(geotags, geotag)
	}

	return geotags, nil
}
