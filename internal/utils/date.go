package utils

import "time"

func AtLeastOneMonthPassed(from, now time.Time) bool {
    loc := from.Location()
    from = time.Date(from.Year(), from.Month(), from.Day(), 12, 0, 0, 0, loc)
    now  = time.Date(now.Year(),  now.Month(),  now.Day(),  12, 0, 0, 0, loc)

    oneMonthLater := from.AddDate(0, 1, 0)
    return !now.Before(oneMonthLater)
}