package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type tripsData struct {
	StartDate string
	Trips     []struct {
		Departure  string
		Return     string
		Comment    string
		IsBusiness bool
	}
}

func readData(filename string) tripsData {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var result tripsData
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func getDate(s string) time.Time {
	result, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return result
}

var (
	allowedYears = flag.Int("allowed_years", 8, "Years allowed for the 4-year period")
	requiredDays = flag.Int("required_days", 1460, "Days required in the 4-year period")
	dataFile     = flag.String("data_file", "data.json", "Data file to read")
)

func main() {
	flag.Parse()

	// TODO(lutzky): Leap years

	data := readData(*dataFile)

	startDate := getDate(data.StartDate)

	i := 0
	inTrip := false
	inBusiness := false
	daysLenient := 1
	daysStrict := 1
	for t := startDate; t.Sub(startDate) < time.Duration(*allowedYears)*365*24*time.Hour; t = t.Add(24 * time.Hour) {
		departureDate := getDate("2999-01-01")
		returnDate := departureDate
		if i < len(data.Trips) {
			departureDate = getDate(data.Trips[i].Departure)
			returnDate = getDate(data.Trips[i].Return)
		}
		header := fmt.Sprintf("%s inTrip:%5t inBusiness:%5t daysLenient:%4d daysStrict:%4d ", t.Format("2006-01-02"), inTrip, inBusiness, daysLenient, daysStrict)
		if daysLenient == *requiredDays {
			fmt.Printf("Reached required days (lenient version) on %s\n", t.Format("2006-01-02"))
		}
		if daysStrict == *requiredDays {
			fmt.Printf("Reached required days (strict version) on %s\n", t.Format("2006-01-02"))
			return
		}
		if !inTrip {
			if t.Before(departureDate.Add(-24 * time.Hour)) {
				fmt.Println(header, "Nothing")
			} else {
				fmt.Println(header, "Started trip: ", data.Trips[i].Comment)
				inTrip = true
				inBusiness = data.Trips[i].IsBusiness
			}
		} else {
			if t.Before(returnDate) {
				fmt.Println(header, "Nothing")

			} else {
				fmt.Println(header, "Ended trip: ", data.Trips[i].Comment)
				inTrip = false
				inBusiness = false
				i++
			}
		}
		if !inTrip {
			daysLenient++
			daysStrict++
		} else {
			if inBusiness {
				daysLenient++
			}
		}

	}

}
