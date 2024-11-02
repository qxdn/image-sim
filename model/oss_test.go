package model_test

import (
	"testing"

	"github.com/qxdn/imagesim/model"
)

func TestExtractFilename(t *testing.T) {
	// Test case 1
	// Test the function with a normal URL
	testExtractFilename(t, "/images/1.jpg", "1.jpg")

	// Test case 2
	// Test the function with a URL that has no filename
	testExtractFilename(t, "/images/", "")

	// Test case 3
	// Test the function with a URL that has a long path
	testExtractFilename(t, "/images/1/2/3/4/5.jpg", "5.jpg")

	// Test case 4
	testExtractFilename(t, "1.jpg", "1.jpg")
}

func testExtractFilename(t *testing.T, input, expect string) {
	filename := model.ExtractOSSFilename(input)
	if filename != expect {
		t.Errorf("Expected filename is %v, but got %v", expect, filename)
	}

}
