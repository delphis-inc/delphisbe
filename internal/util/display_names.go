package util

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"
)

// Pull the entropy from bits 0-9 (most significant, big endian)
func GenerateAnimalDisplayName(seed uint64) string {
	seedShifted := seed >> 54
	return AnimalArray[seedShifted%uint64(len(AnimalArray))]
}

// Pull the entropy from bits 10-15 (most significant, big endian)
func GenerateGradient(seed uint64) model.GradientColor {
	seedShifted := (seed << 10) >> 58
	// NOTE: The first element in the gradient color array is `Unknown` so we
	// do some conversion around that.
	index := (seedShifted % uint64(len(model.AllGradientColor)-1)) + 1

	return model.AllGradientColor[index]
}

// Pull the entropy from bits 16-25 (most significant)
func GenerateDisplayNameIndex(seed uint64) int {
	seedShifted := (seed << 16) >> 54

	return int(seedShifted % uint64(1024))
}

func GenerateFullDisplayName(seed uint64) string {
	return strings.Title(fmt.Sprintf("%s %s (#%d)",
		strings.ToLower(string(GenerateGradient(seed))),
		strings.ToLower(GenerateAnimalDisplayName(seed)),
		GenerateDisplayNameIndex(seed)),
	)
}

func GenerateParticipantSeed(discussionID, participantID string, shuffleCount int) uint64 {
	// We generate the display name by SHA-1(discussion_id, participant.id, shuffle_count) without
	// commas, just concatenated.
	h := sha1.Sum([]byte(fmt.Sprintf("%s%s%d", discussionID, participantID, shuffleCount)))
	return binary.BigEndian.Uint64(h[:])
}
