package blia

var defaultPageLimit = 20
var maxPageLimit = 999

func SetMaxPageLimit(n int) {
	maxPageLimit = n
}

func SetDefaultPageLimit(n int) {
	defaultPageLimit = n
}
