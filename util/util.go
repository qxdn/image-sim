package util

import (
	"image"

	"math/bits"

	"github.com/corona10/goimagehash"
)

/**
 * @Function: AHash
 * @Description: Calculate the average hash value of the image.
 * @Param: img image.Image
 * @Return: uint64, error
 **/
func AHash(img *image.Image) (uint64, error) {
	ahash, err := goimagehash.AverageHash(*img)
	if err != nil {
		return 0, err
	}
	return ahash.GetHash(), nil
}

/**
 * @Function: DHash
 * @Description: Calculate the difference hash value of the image.
 * @Param: img image.Image
 * @Return: uint64, error
 **/
func DHash(img *image.Image) (uint64, error) {
	dhash, err := goimagehash.DifferenceHash(*img)
	if err != nil {
		return 0, err
	}
	return dhash.GetHash(), nil
}

/**
 * @Function: PHash
 * @Description: Calculate the perceptual hash value of the image.
 * @Param: img image.Image
 * @Return: uint64, error
 **/
func PHash(img *image.Image) (uint64, error) {
	phash, err := goimagehash.PerceptionHash(*img)
	if err != nil {
		return 0, err
	}
	return phash.GetHash(), nil
}

func HashDistance(hash1, hash2 uint64) int {
	// 汉明距离
	hamming := hash1 ^ hash2
	return bits.OnesCount64(hamming)
}
