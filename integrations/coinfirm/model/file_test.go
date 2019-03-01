package model_test

import (
	"testing"

	"modulus/kyc/integrations/coinfirm/model"

	"github.com/stretchr/testify/assert"
)

func TestIsAcceptedFileExt(t *testing.T) {
	assert := assert.New(t)

	ext := "jpg"
	isAccepted := model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "jpeg"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "png"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "gif"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "bmp"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "svg"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "psd"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "tif"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "tiff"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "webp"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	ext = "pdf"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.True(isAccepted)

	// Unknown extensions should fail.

	ext = "docx"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "pptx"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "cdr"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	// Uppercase or mixed case extensions should fail.

	ext = "JPG"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "JpeG"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = "Png"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	// Non-normalized extensions should fail.

	ext = ".jpeg"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)

	ext = ".png"
	isAccepted = model.IsAcceptedFileExt(ext)

	assert.False(isAccepted)
}
