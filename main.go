package main

import (
        "encoding/json"
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

const allowedYears = 8
const requiredDays = 1460

func main() {
        // TODO(lutzky): Leap years

        data := readData(os.Args[1])

        startDate := getDate(data.StartDate)

        i := 0
        inTrip := false
        inBusiness := false
        daysLenient := 0
        daysStrict := 0
        for t := startDate; t.Sub(startDate) < allowedYears*365*24*time.Hour; t = t.Add(24 * time.Hour) {
                if daysLenient == requiredDays {
                        fmt.Printf("Reached required days (lenient version) on %s\n", t.Format("2006-01-02"))
                }
                if daysStrict == requiredDays {
                        fmt.Printf("Reached required days (strict version) on %s\n", t.Format("2006-01-02"))
                        return
                }
                departureDate := getDate("2999-01-01")
                returnDate := departureDate
                if i < len(data.Trips) {
                        departureDate = getDate(data.Trips[i].Departure)
                        returnDate = getDate(data.Trips[i].Return)
                }
                if !inTrip {
                        daysLenient++
                        daysStrict++
                } else {
                        if inBusiness {
                                daysLenient++
                        }
                }
                header := fmt.Sprintf("%s inTrip:%5t inBusiness:%5t daysLenient:%4d daysStrict:%4d ", t.Format("2006-01-02"), inTrip, inBusiness, daysLenient, daysStrict)
                if !inTrip {
                        if t.Before(departureDate) {
                                //fmt.Println(header, "Nothing")
                                continue
                        }
                        fmt.Println(header, "Started trip: ", data.Trips[i].Comment)
                        inTrip = true
                        inBusiness = data.Trips[i].IsBusiness
                } else {
                        if t.Before(returnDate) {
                                //fmt.Println(header, "Nothing")
                                continue
                        }
                        fmt.Println(header, "Ended trip: ", data.Trips[i].Comment)
                        inTrip = false
                        inBusiness = false
                        i++
                }
        }

}

