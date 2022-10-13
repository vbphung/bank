package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomFullName() string {
	return fmt.Sprintf("%s %s", pickLastName(), pickFirstName())
}

func RandomBalance(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func pickLastName() string {
	lastNames := []string{"Hightower", "Stark", "Lannister", "Targaryen", "Velaryon"}
	return pickName(lastNames)
}

func pickFirstName() string {
	firstNames := []string{"Alicent", "Otto", "Corlys", "Jon", "Jaime"}
	return pickName(firstNames)
}

func pickName(names []string) string {
	return names[rand.Intn(len(names))]
}
