package main

import "strconv"

var counter uint64
var ids map[string]string
var venue_type map[string]VenueType

func init() {
	counter = 0
	ids = make(map[string]string)
	venue_type = make(map[string]VenueType)
}

// getId returns the ID of the given value. If the value does not exist, it is
// assigned from counter and the counter is increased
func getId(value string) (string, bool) {
	if val, ok := ids[value]; ok {
		return val, true
	}

	counter++
	scount := strconv.FormatUint(counter, 10)
	ids[value] = scount

	return scount, false
}
