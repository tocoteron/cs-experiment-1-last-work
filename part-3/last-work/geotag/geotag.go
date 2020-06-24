package geotag

import (
	"encoding/csv"
	"math/rand"
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
		URL:       data[3],
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

func LoadGeoTags(path string) ([]GeoTag, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	geotags, err := UnmarshalGeoTags(lines)
	if err != nil {
		return nil, err
	}

	return geotags, nil
}

func GenerateRandomGeoTag(id int, tagLen int) GeoTag {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, tagLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	geotag := GeoTag{
		ID:        id,
		Time:      string(b),
		Latitude:  rand.Float64(),
		Longitude: rand.Float64(),
		URL:       string(b),
	}

	return geotag
}

func GenerateRandomGeoTags(num int, tagLen int) []GeoTag {
	geotags := make([]GeoTag, num)

	for i := 0; i < num; i++ {
		geotags[i] = GenerateRandomGeoTag(i, tagLen)
	}

	return geotags
}
