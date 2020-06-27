package geotag

import (
	"cs-experiment-1/part-3/last-work/csvutil"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// GeoTag is corresponding the db table 'geotag'
type GeoTag struct {
	ID        uint64
	Time      int32
	Latitude  float64
	Longitude float64
	URLFarmID uint8
	URLID1    uint64
	URLID2    uint64
}

// IDSearchTable is mapping ID to GeoTag
type IDSearchTable map[uint64]*GeoTag

// TagSearchTable is mapping Tag to Geotag
type TagSearchTable map[string][]*GeoTag

// Datetime is raw datetime info
func (geotag *GeoTag) Datetime() string {
	t := time.Unix(int64(geotag.Time), 0)
	return t.Format("2006-01-02 15:04:05")
}

// URL is raw url info
func (geotag *GeoTag) URL() string {
	return fmt.Sprintf(
		"http://farm%d.static.flickr.com/%d/%d_%010x.jpg",
		geotag.URLFarmID,
		geotag.URLID1,
		geotag.ID,
		geotag.URLID2,
	)
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

	var id uint64
	var urlFarmID uint8
	var urlID1 uint64
	var urlID2 uint64

	_, err = fmt.Sscanf(
		data[4],
		"http://farm%d.static.flickr.com/%d/%d_%x.jpg",
		&urlFarmID,
		&urlID1,
		&id,
		&urlID2,
	)
	if err != nil {
		return GeoTag{}, err
	}

	geotag := GeoTag{
		ID:        id,
		Time:      datetime,
		Latitude:  latitude,
		Longitude: longitude,
		URLFarmID: urlFarmID,
		URLID1:    urlID1,
		URLID2:    urlID2,
	}

	return geotag, nil
}

func UnmarshalCompressedGeoTag(data []string) (GeoTag, error) {
	id, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return GeoTag{}, err
	}

	datetime, err := strconv.ParseInt(data[1], 10, 32)
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

	urlFarmID, err := strconv.ParseUint(data[4], 10, 8)
	if err != nil {
		return GeoTag{}, err
	}

	urlID1, err := strconv.ParseUint(data[5], 10, 64)
	if err != nil {
		return GeoTag{}, err
	}

	urlID2, err := strconv.ParseUint(data[6], 10, 64)
	if err != nil {
		return GeoTag{}, err
	}

	geotag := GeoTag{
		ID:        id,
		Time:      int32(datetime),
		Latitude:  latitude,
		Longitude: longitude,
		URLFarmID: uint8(urlFarmID),
		URLID1:    urlID1,
		URLID2:    urlID2,
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

func ReadCompressedGeoTagsFromCSV(path string, capacity int, buffsize int) ([]GeoTag, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	geotags := make([]GeoTag, 0, capacity)

	for record := range csvutil.AsyncReadCSV(reader, buffsize) {
		geotag, err := UnmarshalCompressedGeoTag(record)
		if err != nil {
			return nil, err
		}

		geotags = append(geotags, geotag)
	}

	return geotags, nil
}

func WriteGeoTagsToCSV(path string, geotags []GeoTag) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	records := make([][]string, len(geotags))

	for i := 0; i < len(geotags); i++ {
		records[i] = []string{
			strconv.FormatUint(geotags[i].ID, 10),
			strconv.Itoa(int(geotags[i].Time)),
			strconv.FormatFloat(geotags[i].Latitude, 'f', -1, 64),
			strconv.FormatFloat(geotags[i].Longitude, 'f', -1, 64),
			strconv.FormatUint(uint64(geotags[i].URLFarmID), 10),
			strconv.FormatUint(uint64(geotags[i].URLID1), 10),
			strconv.FormatUint(uint64(geotags[i].URLID2), 10),
		}
	}

	writer := csv.NewWriter(file)
	writer.WriteAll(records)
	writer.Flush()

	return nil
}

func WriteTagSearchTableToCSV(path string, tagSearchTable TagSearchTable) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	records := make([][]string, len(tagSearchTable))

	i := 0
	for k, geotagPointers := range tagSearchTable {
		records[i] = []string{k}

		for _, geotagPointer := range geotagPointers {
			records[i] = append(records[i], strconv.FormatUint(geotagPointer.ID, 10))
		}

		i++
	}

	fmt.Println(len(records[0]))

	writer := csv.NewWriter(file)
	writer.WriteAll(records)
	writer.Flush()

	return nil
}

func ReadTagSearchTableFromCSV(path string, idSearchTable IDSearchTable) (TagSearchTable, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tagSearchTable := make(TagSearchTable)

	for record := range csvutil.AsyncReadCSV(reader, 1000) {
		tag := record[0]

		for i := 1; i < len(record); i++ {
			id, err := strconv.ParseUint(record[i], 10, 64)
			if err != nil {
				return TagSearchTable{}, err
			}

			tagSearchTable[tag] = append(tagSearchTable[tag], idSearchTable[id])
		}
	}

	return tagSearchTable, nil
}
