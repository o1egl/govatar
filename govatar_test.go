package govatar

import (
	"image"
	"image/color"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	testCases := map[string]struct {
		gender Gender
		err    error
	}{
		"male": {
			gender: MALE,
		},
		"female": {
			gender: FEMALE,
		},
		"unsupported": {
			gender: Gender(-1),
			err:    ErrUnsupportedGender,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			avatar, err := Generate(test.gender)
			assert.Equal(t, test.err, err)
			if test.err == nil && assert.NotNil(t, avatar) {
				bounds := avatar.Bounds()
				if assert.False(t, bounds.Empty()) {
					assert.Equal(t, 400, bounds.Dx())
					assert.Equal(t, 400, bounds.Dy())
				}
			}
		})
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestGenerateFile(t *testing.T) {
	fileName := randStringRunes(6)
	testCases := map[string]struct {
		imageName string
		mimeType  string
	}{
		"jpeg": {imageName: fileName + ".jpeg", mimeType: "image/jpeg"},
		"jpg":  {imageName: fileName + ".jpg", mimeType: "image/jpeg"},
		"png":  {imageName: fileName + ".png", mimeType: "image/png"},
		"gif":  {imageName: fileName + ".gif", mimeType: "image/gif"},
		"xyz":  {imageName: fileName + ".xyz", mimeType: "image/png"},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			os.Remove(test.imageName)
			defer os.Remove(test.imageName)
			err := GenerateFile(MALE, test.imageName)
			if assert.NoError(t, err) {
				buf := make([]byte, 512)
				f, err := os.Open(test.imageName)
				if assert.NoError(t, err) {
					defer f.Close()

					_, err = f.Read(buf)
					if assert.NoError(t, err) {
						assert.Equal(t, test.mimeType, http.DetectContentType(buf))
					}
				}
			}
		})
	}
}

func TestGenerateFromString(t *testing.T) {
	testCases := map[string]struct {
		gender    Gender
		areEqual  bool
		username1 string
		username2 string
		err       error
	}{
		"equal": {
			gender:    MALE,
			areEqual:  true,
			username1: "username@site.com",
			username2: "username@site.com",
		},
		"not equal": {
			gender:    MALE,
			areEqual:  false,
			username1: "username@site.com",
			username2: "another-username@site.com",
		},
		"error": {
			gender:    Gender(-1),
			username1: "username@site.com",
			username2: "username@site.com",
			err:       ErrUnsupportedGender,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			avatar1, err := GenerateForUsername(test.gender, test.username1)
			assert.Equal(t, test.err, err)
			if test.err == nil && assert.NotNil(t, avatar1) {
				bounds := avatar1.Bounds()
				assert.False(t, bounds.Empty())
				assert.Equal(t, 400, bounds.Dx())
				assert.Equal(t, 400, bounds.Dy())

				avatar2, err := GenerateForUsername(test.gender, test.username2)
				if assert.NoError(t, err) {
					assert.Equal(t, test.areEqual, areImagesEqual(avatar1, avatar2))
				}
			}
		})
	}

}

func TestGenerateFileFromString(t *testing.T) {
	fileName := randStringRunes(6)
	testCases := map[string]struct {
		imageName string
		mimeType  string
	}{
		"jpeg": {imageName: fileName + ".jpeg", mimeType: "image/jpeg"},
		"jpg":  {imageName: fileName + ".jpg", mimeType: "image/jpeg"},
		"png":  {imageName: fileName + ".png", mimeType: "image/png"},
		"gif":  {imageName: fileName + ".gif", mimeType: "image/gif"},
		"xyz":  {imageName: fileName + ".xyz", mimeType: "image/png"},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			os.Remove(test.imageName)
			defer os.Remove(test.imageName)
			err := GenerateFileForUsername(MALE, "username@site.com", test.imageName)
			if assert.NoError(t, err) {
				buf := make([]byte, 512)
				f, err := os.Open(test.imageName)
				if assert.NoError(t, err) {
					defer f.Close()

					_, err = f.Read(buf)
					if assert.NoError(t, err) {
						assert.Equal(t, test.mimeType, http.DetectContentType(buf))
					}
				}
			}
		})
	}
}

func areImagesEqual(a, b image.Image) bool {
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
