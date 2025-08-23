package livetiming

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := os.Open("./../../data/2021/2021-04-18_Emilia_Romagna_Grand_Prix/2021-04-18_Race/CarData.z.jsonStream")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = Parse(file)
	if err != nil {
		t.Error(err)
	}
}
