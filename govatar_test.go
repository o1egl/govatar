package govatar

import (
	"image"
	"image/color"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	// Male test
	avatar, err := Generate(MALE)
	bounds := avatar.Bounds()

	assert.NoError(t, err)
	assert.NotNil(t, avatar)
	assert.False(t, bounds.Empty())
	assert.Equal(t, 400, bounds.Dx())
	assert.Equal(t, 400, bounds.Dy())

	// Female test
	avatar, err = Generate(FEMALE)
	bounds = avatar.Bounds()

	assert.NoError(t, err)
	assert.NotNil(t, avatar)
	assert.False(t, bounds.Empty())
	assert.Equal(t, 400, bounds.Dx())
	assert.Equal(t, 400, bounds.Dy())
}

func TestGenerateFile(t *testing.T) {
	generateFileTest(t, MALE)
	generateFileTest(t, FEMALE)
}

func generateFileTest(t *testing.T, gender Gender) {
	imagesTests := []struct {
		imageName string
		mimeType  string
	}{
		{"avatar.jpeg", "image/jpeg"},
		{"avatar.jpg", "image/jpeg"},
		{"avatar.png", "image/png"},
		{"avatar.gif", "image/gif"},
		{"avatar.xyz", "image/png"},
	}
	// test jpeg generation
	for _, tt := range imagesTests {
		os.Remove(tt.imageName)
		err := GenerateFile(gender, tt.imageName)
		assert.Nil(t, err)

		buf := make([]byte, 512)
		f, err := os.Open(tt.imageName)
		assert.Nil(t, err)
		defer f.Close()

		_, err = f.Read(buf)
		assert.Nil(t, err)

		assert.Equal(t, tt.mimeType, http.DetectContentType(buf))
	}
}

func TestGenerateFromString(t *testing.T) {
	// Male test
	avatar1, err := GenerateForUsername(MALE, "username@site.com")
	bounds := avatar1.Bounds()

	assert.NoError(t, err)
	assert.NotNil(t, avatar1)
	assert.False(t, bounds.Empty())
	assert.Equal(t, 400, bounds.Dx())
	assert.Equal(t, 400, bounds.Dy())

	avatar2, err := GenerateForUsername(MALE, "username@site.com")
	assert.NoError(t, err)
	assert.True(t, areImagesEquals(avatar1, avatar2))

	// Female test
	avatar1, err = GenerateForUsername(FEMALE, "username@site.com")
	assert.NoError(t, err)

	avatar2, err = GenerateForUsername(FEMALE, "username2@site.com")
	assert.NoError(t, err)

	assert.False(t, areImagesEquals(avatar1, avatar2))

}

func TestGenerateFileFromString(t *testing.T) {
	generateFileFromStringTest(t, MALE)
	generateFileFromStringTest(t, FEMALE)
}

func generateFileFromStringTest(t *testing.T, gender Gender) {
	imagesTests := []struct {
		imageName string
		mimeType  string
	}{
		{"avatar.jpeg", "image/jpeg"},
		{"avatar.jpg", "image/jpeg"},
		{"avatar.png", "image/png"},
		{"avatar.gif", "image/gif"},
		{"avatar.xyz", "image/png"},
	}

	for _, tt := range imagesTests {
		os.Remove(tt.imageName)
		err := GenerateFileForUsername(gender, "username@site.com", tt.imageName)
		assert.Nil(t, err)

		buf := make([]byte, 512)
		f, err := os.Open(tt.imageName)
		assert.Nil(t, err)
		defer f.Close()

		_, err = f.Read(buf)
		assert.Nil(t, err)

		assert.Equal(t, tt.mimeType, http.DetectContentType(buf))
	}
}

func areImagesEquals(a, b image.Image) bool {
	ab, bb := a.Bounds(), b.Bounds()
	w, h := ab.Dx(), ab.Dy()
	if w != bb.Dx() || h != bb.Dy() {
		return false
	}
	n := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d := diffColor(a.At(ab.Min.X+x, ab.Min.Y+y), b.At(bb.Min.X+x, bb.Min.Y+y))
			c := color.RGBA{0, 0, 0, 0xff}
			if d > 0 {
				c.R = 0xff
				n++
			}
		}
	}
	return n == 0
}

func diffColor(c1, c2 color.Color) int64 {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	var diff int64
	diff += abs(int64(r1) - int64(r2))
	diff += abs(int64(g1) - int64(g2))
	diff += abs(int64(b1) - int64(b2))
	diff += abs(int64(a1) - int64(a2))
	return diff
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
