package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FundStatus struct {
	PortfolioPercent float32
	Date             time.Time
	CustomDate       string
}

func main() {
	filesPaths, err := GetFilePaths()
	if err != nil {
		return
	}
	if len(filesPaths) == 0 {
		log.Println("no files in drop-files-here folder")
		return
	}

	for _, filepath := range filesPaths {
		strFileContentArr, err := DecodeFileToArray(filepath)
		if err != nil {
			return
		}

		fundArr, err := SeparateElementsIntoStructArray(strFileContentArr)
		if err != nil {
			log.Println(err)
			return
		}

		// ConvertStructToJson(fundArr)
		ProcessData(fundArr)
	}
}

func GetFilePaths() (filesPaths []string, err error) {
	err = filepath.Walk("drop-files-here", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println("walk function error", err)
			return err
		}

		if !info.IsDir() {
			filesPaths = append(filesPaths, path)
		}

		return nil
	})
	if err != nil {
		log.Println("error occurred", err)
		return
	}
	return
}

func DecodeFileToArray(filepath string) (strFileContentArr []string, err error) {
	byteFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Println("error occurred while reading file", err)
		return []string{}, err
	}

	stringFile := string(byteFile)

	strFileContentArr = strings.Fields(stringFile)

	return
}

func SeparateElementsIntoStructArray(strarr []string) (fundarr []FundStatus, err error) {
	fundarr = make([]FundStatus, 0)

	i, j, k, l := 0, 1, 2, 3
	for i < len(strarr)-1 {

		strpercent := strarr[i]
		strpercent = strings.Trim(strpercent, "\ufeff")
		var floatPercent float64
		floatPercent, err = strconv.ParseFloat(strpercent, 32)
		if err != nil {
			log.Println("float conversion error", err)
			return nil, err
		}

		percent := float32(floatPercent)

		var day int
		day, err = strconv.Atoi(strarr[j])
		if err != nil {
			log.Println("int conversion error", err)
			return nil, err
		}

		month := strarr[k]
		timeMonth := ExtractMonth(month)

		var year int
		year, err = strconv.Atoi(strarr[l])
		if err != nil {
			log.Println("int conversion error", err)
			return nil, err
		}

		datetime := time.Date(year, timeMonth, day, 0, 0, 0, 0, time.UTC)

		customDateString := datetime.Format("02-Jan-06")

		singlefundstatus := FundStatus{
			PortfolioPercent: percent,
			Date:             datetime,
			CustomDate:       customDateString,
		}
		fundarr = append(fundarr, singlefundstatus)

		i, j, k, l = i+4, j+4, k+4, l+4
		if i > len(strarr) || j > len(strarr) || k > len(strarr) {
			break
		}
	}

	return
}

func ProcessData(fundarr []FundStatus) {
	DataCalculationByDays(fundarr)
}

func DataCalculationByDays(fundarr []FundStatus) {
	var allDaysAverage float32

	for _, singleRecord := range fundarr {
		allDaysAverage = allDaysAverage + singleRecord.PortfolioPercent
	}

	firstDay := fundarr[0].Date
	lastDay := fundarr[len(fundarr)-1].Date
	duration := lastDay.Sub(firstDay)
	totalNumberOfDays := float32(duration.Hours() / 24)

	var avg1, avg2 float32
	if len(fundarr) < int(totalNumberOfDays) {
		avg1 = allDaysAverage / float32(len(fundarr))
		log.Println("complete portfolio average for", len(fundarr), "days between", fundarr[0].CustomDate, "to", fundarr[len(fundarr)-1].CustomDate, "is", avg1)

	} else {
		avg2 = allDaysAverage / totalNumberOfDays
		log.Println("complete portfolio average for", totalNumberOfDays, "days between", fundarr[0].CustomDate, "to", fundarr[len(fundarr)-1].CustomDate, "is", avg2)
	}
}

func GetTotalDaysOfMonth(year int) (numberOfDaysInMonth int) {
	month := time.September

	firstDayOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-time.Second)

	numberOfDaysInMonth = lastDayOfCurrentMonth.Day()

	fmt.Printf("Number of days in %s %d: %d\n", month, year, numberOfDaysInMonth)
	return
}

func ExtractMonth(monthStr string) (month time.Month) {
	switch monthStr {
	case "jan":
		month = time.January
	case "feb":
		month = time.February
	case "mar":
		month = time.March
	case "apr":
		month = time.April
	case "may":
		month = time.May
	case "jun":
		month = time.June
	case "jul":
		month = time.July
	case "aug":
		month = time.August
	case "sep":
		month = time.September
	case "oct":
		month = time.October
	case "nov":
		month = time.November
	case "dec":
		month = time.December
	}

	return
}

func ConvertStructToJson(fundarr []FundStatus) {
	bytefunds, err := json.MarshalIndent(fundarr, " ", "  ")
	if err != nil {
		log.Println("json converson error", err)
		return
	}

	strfunds := string(bytefunds)

	log.Printf("%+v", strfunds)
}
