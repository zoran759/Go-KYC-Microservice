package model_test

import (
	"testing"

	"modulus/kyc/integrations/coinfirm/model"

	"github.com/stretchr/testify/assert"
)

func TestIsAcceptedFileExt(t *testing.T) {
	assert := assert.New(t)

	ext := "doc"
	isAccepted := model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "psd"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "TIF"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = model.NormalizeFileExt(ext)
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "jpeg"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = model.NormalizeFileExt(ext)
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "JPEG"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = model.NormalizeFileExt(ext)
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "jpg"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "PNG"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = model.NormalizeFileExt(ext)
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "Pdf"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = model.NormalizeFileExt(ext)
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)
}

func TestNormalizeFileExt(t *testing.T) {
	assert := assert.New(t)

	ext := "doc"
	next := model.NormalizeFileExt(ext)

	assert.Equal(ext, next)

	ext2 := "Doc"
	next = model.NormalizeFileExt(ext2)

	assert.Equal(ext, next)

	ext2 = "DOC"
	next = model.NormalizeFileExt(ext2)

	assert.Equal(ext, next)

	ext = "jpg"
	ext2 = "jpeg"
	next = model.NormalizeFileExt(ext2)

	assert.Equal(ext, next)

	ext2 = "JPEG"
	next = model.NormalizeFileExt(ext2)

	assert.Equal(ext, next)

	ext2 = "JPeg"
	next = model.NormalizeFileExt(ext2)

	assert.Equal(ext, next)
}
