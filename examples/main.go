package main

import (
	"fmt"
	"time"

	"github.com/gahojin/go-holiday-japanese"
)

const (
	DateLayout = "2006-01-02"
)

func main() {
	t, _ := time.Parse(DateLayout, "2024-10-14")
	fmt.Printf("%s = %s\n", t, holiday.GetHolidayName(t))

	t, _ = time.Parse(DateLayout, "2024-10-15")
	fmt.Printf("%s = %s\n", t, holiday.GetHolidayName(t))

	s, _ := time.Parse(DateLayout, "2025-01-01")
	e, _ := time.Parse(DateLayout, "2025-12-31")
	fmt.Printf("%s ~ %s = %v\n", s, e, holiday.Between(s, e))
}
