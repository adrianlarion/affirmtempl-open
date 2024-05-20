package csv

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adrianlarion/affirmtempl-open/internal/model"
)

func DefaultCardsFromCsv() ([]*model.AffirmCard, error) {
	affirmCards := []*model.AffirmCard{}
	csvDir := "./data/csv/cards"

	csvFiles, err := csvFilesInDir(csvDir)
	if err != nil {
		return nil, err
	}
	for _, csvFile := range csvFiles {
		c, err := csvToAffirmCard(csvFile)
		if err != nil {
			return nil, err
		}
		affirmCards = append(affirmCards, c)

	}

	return affirmCards, nil
}

func csvFilesInDir(csvDir string) ([]string, error) {
	var csvFiles []string
	files, err := os.ReadDir(csvDir)

	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {

			fp := filepath.Join(csvDir, file.Name())
			csvFiles = append(csvFiles, fp)
		}
	}
	return csvFiles, nil
}

func csvToRecords(csvPath string) ([][]string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("error while reading csv", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("error while reading recors", err)
		return nil, err
	}
	// for _, eachrecord := range records{
	// 	fmt.Println(eachrecord)
	// }

	return records, nil
}

func recordsAreCorrect(records [][]string) bool {
	if len(records) < 3 {
		return false
	}
	return true
}

func csvToAffirmCard(csvPath string) (*model.AffirmCard, error) {
	records, err := csvToRecords(csvPath)
	if err != nil {
		return nil, err
	}
	if !recordsAreCorrect(records) {
		return nil, errors.New("records from csv are not correcct")
	}

	// fmt.Println(records)
	var a = &model.AffirmCard{}
	a.Title = records[0][0]

	fav, err := strconv.ParseBool(records[1][0])
	if err != nil {
		return nil, errors.New("favorite record from csv cannot be parsed")
	}
	a.Fav = fav
	a.ImgPath = records[2][0]

	id, err := strconv.ParseInt(records[3][0], 10, 64)
	if err != nil {
		return nil, errors.New("id record from csv cannot be parsed")
	}
	a.ID = id

	a.DefaultAffirmations = []model.AffirmEntry{}
	if len(records) >= 5 {

		for i := 4; i < len(records); i++ {

			a.DefaultAffirmations = append(a.DefaultAffirmations, model.AffirmEntry{Content: records[i][0]})

		}
	}
	//copy default to affirmations
	a.Affirmations = make([]model.AffirmEntry, len(a.DefaultAffirmations))
	copy(a.Affirmations, a.DefaultAffirmations)

	return a, nil
}
