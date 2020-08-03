package util

import (
	"fmt"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"
)

// Pull the entropy from the 8 most significant bits
func GenerateAnimalDisplayName(seed uint64) string {
	seedShifted := seed >> 56
	return AnimalArray[seedShifted%uint64(len(AnimalArray))]
}

// Pull the entropy from the second 8 most significant bits
func GenerateGradient(seed uint64) model.GradientColor {
	seedShifted := (seed << 8) >> 56
	// NOTE: The first element in the gradient color array is `Unknown` so we
	// do some conversion around that.
	index := (seedShifted % uint64(len(model.AllGradientColor)-1)) + 1

	return model.AllGradientColor[index]
}

// Pull the entropy from the third 8 most significant bits
func GenerateDisplayNameIndex(seed uint64) int {
	seedShifted := (seed << 16) >> 56

	return int(seedShifted % uint64(1024))
}

func GenerateFullDisplayName(seed uint64) string {
	return strings.Title(fmt.Sprintf("%s %s (#%d)",
		strings.ToLower(string(GenerateGradient(seed))),
		strings.ToLower(GenerateAnimalDisplayName(seed)),
		GenerateDisplayNameIndex(seed)),
	)
}
