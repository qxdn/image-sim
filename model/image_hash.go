package model

import (
	"image"

	"github.com/qxdn/imagesim/util"
)

type ImageHash struct {
	AHash uint64 // Average hash
	DHash uint64 // Difference hash
	PHash uint64 // Perception hash
}

/**
 * @Function: ComputeHashes
 * @Description: Compute all hashes of the image.
 * @Param: img *image.Image
 * @Return: *ImageHash, error
 **/
func ComputeHashes(img *image.Image) (*ImageHash, error) {
	hash := &ImageHash{}
	err := hash.ComputeHashes(img)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

/**
 * @Function: ComputeHashes
 * @Description: Compute all hashes of the image.
 * @Param: img *image.Image
 * @Return: error
 **/
func (i *ImageHash) ComputeHashes(img *image.Image) error {
	_, err := i.ComputeAHash(img)
	if err != nil {
		return err
	}

	_, err = i.ComputeDHash(img)
	if err != nil {
		return err
	}

	_, err = i.ComputePHash(img)
	if err != nil {
		return err
	}

	return nil
}

func (i *ImageHash) ComputeAHash(img *image.Image) (uint64, error) {
	ahash, err := util.AHash(img)
	if err != nil {
		return 0, err
	}
	i.AHash = ahash
	return ahash, nil
}

func (i *ImageHash) ComputeDHash(img *image.Image) (uint64, error) {
	dhash, err := util.DHash(img)
	if err != nil {
		return 0, err
	}
	i.DHash = dhash
	return dhash, nil
}

func (i *ImageHash) ComputePHash(img *image.Image) (uint64, error) {
	phash, err := util.PHash(img)
	if err != nil {
		return 0, err
	}
	i.PHash = phash
	return phash, nil
}

func (i *ImageHash) GetAHash() uint64 {
	return i.AHash
}

func (i *ImageHash) GetDHash() uint64 {
	return i.DHash
}

func (i *ImageHash) GetPHash() uint64 {
	return i.PHash
}

func (i *ImageHash) AHashDistance(other *ImageHash) int {
	return util.HashDistance(i.AHash, other.AHash)
}

func (i *ImageHash) DHashDistance(other *ImageHash) int {
	return util.HashDistance(i.DHash, other.DHash)
}

func (i *ImageHash) PHashDistance(other *ImageHash) int {
	return util.HashDistance(i.PHash, other.PHash)
}
