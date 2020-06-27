package csvutil

import (
	"encoding/csv"
	"io"
	"log"
)

func AsyncReadCSV(ioreader io.Reader, buffsize int) chan []string {
	reader := csv.NewReader(ioreader)
	ch := make(chan []string, buffsize)

	go func() {
		defer close(ch)
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			ch <- record
		}
	}()

	return ch
}
