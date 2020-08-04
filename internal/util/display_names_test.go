package util

import (
	"fmt"
	"strings"
	"testing"

	"github.com/delphis-inc/delphisbe/graph/model"
	. "github.com/smartystreets/goconvey/convey"
)

//zero := uint64(0)
var allOnes = ^uint64(0)

func TestUtils_GenerateAnimalDisplayName(t *testing.T) {
	Convey("GenerateAnimalDisplayName", t, func() {
		Convey("when 0 in first 10 bits (big endian) return 0th element", func() {
			zeroesInBegin := allOnes >> 10

			animalDisplayName := GenerateAnimalDisplayName(zeroesInBegin)

			So(animalDisplayName, ShouldEqual, AnimalArray[0])
		})
		Convey("when 1023 in first 10 bits, return 1023 % num animals", func() {
			animalDisplayName := GenerateAnimalDisplayName(allOnes)

			So(animalDisplayName, ShouldEqual, AnimalArray[1023%len(AnimalArray)])
		})
	})
}

func TestUtils_GenerateGradient(t *testing.T) {
	Convey("GenerateGradient", t, func() {
		Convey("when 0s in bits 10-15 returns 1st element of array (0th is unknown)", func() {
			sixBits := uint64(63) << 48
			onlySixBits := sixBits ^ allOnes

			gradient := GenerateGradient(onlySixBits)

			So(gradient, ShouldEqual, model.AllGradientColor[1])
		})

		Convey("when bits 10-15 return 1, send second element of array", func() {
			sixBitsMinus1 := uint64(62) << 48
			onlySixBits := sixBitsMinus1 ^ allOnes

			gradient := GenerateGradient(onlySixBits)

			So(gradient, ShouldEqual, model.AllGradientColor[2])
		})
	})
}

func TestUtils_GenerateDisplaynameIndex(t *testing.T) {
	Convey("DisplayNameIndex", t, func() {
		Convey("Gathers bits from 16-25 bit", func() {
			tenBits := uint64(1023) << 38
			onlyTenBits := tenBits ^ allOnes

			idx := GenerateDisplayNameIndex(onlyTenBits)

			So(idx, ShouldEqual, 0)
		})
	})
}

func TestUtils_GenerateFullDisplayName(t *testing.T) {
	Convey("GenerateFullDisplayName", t, func() {
		Convey("when all zeroes", func() {
			zeroesOnLeft := allOnes >> 26

			displayName := GenerateFullDisplayName(zeroesOnLeft)

			So(displayName, ShouldEqual,
				strings.Title(fmt.Sprintf("%s %s (#%d)",
					strings.ToLower(string(GenerateGradient(zeroesOnLeft))),
					strings.ToLower(string(GenerateAnimalDisplayName(zeroesOnLeft))),
					GenerateDisplayNameIndex(zeroesOnLeft))))
		})
	})
}
