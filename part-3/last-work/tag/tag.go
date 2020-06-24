package tag

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
)

// Tag is correspondign the db table 'tag'
type Tag struct {
	ID  int
	Tag string
}

func UnmarshalTag(data []string) (Tag, error) {
	id, err := strconv.Atoi(data[0])
	if err != nil {
		return Tag{}, err
	}

	tag := Tag{
		ID:  id,
		Tag: data[1],
	}

	return tag, nil
}

func UnmarshalTags(data [][]string) ([]Tag, error) {
	var tags []Tag

	for i := 0; i < len(data); i++ {
		tag, err := UnmarshalTag(data[i])
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func ReadTagsFromCSV(path string) ([]Tag, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.Read()
	if err != nil {
		return nil, err
	}

	tags, err := UnmarshalTags([][]string{lines})
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func GenerateRandomTag(id int, tagLen int) Tag {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, tagLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	tag := Tag{
		ID:  id,
		Tag: string(b),
	}

	return tag
}

func GenerateRandomTags(num int, tagLen int) []Tag {
	tags := make([]Tag, num)

	for i := 0; i < num; i++ {
		tags[i] = GenerateRandomTag(i, tagLen)
	}

	return tags
}

func WriteTagsToCSV(path string, tags []Tag) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines [][]string

	for i := 0; i < len(tags); i++ {
		lines = append(lines, []string{strconv.Itoa(tags[i].ID), tags[i].Tag})
	}

	writer := csv.NewWriter(file)
	writer.WriteAll(lines)
}
