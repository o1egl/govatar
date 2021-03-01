package govatar

import (
	"bytes"
	"embed"
	"errors"
	"hash/fnv"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var ErrUnsupportedGender = errors.New("unsupported gender")

type person struct {
	Clothes []string
	Eye     []string
	Face    []string
	Hair    []string
	Mouth   []string
}

type store struct {
	Background []string
	Male       person
	Female     person
}

var assetsStore *store

// Gender represents gender type
type Gender int

// Male and female constants
const (
	MALE Gender = iota
	FEMALE
)

//go:embed data/*
var dataFS embed.FS

func init() {
	male := getPerson(MALE)
	female := getPerson(FEMALE)
	assetsStore = &store{Background: mustAssetsList("data/background"), Male: male, Female: female}
	rand.Seed(time.Now().UTC().UnixNano())
}

// Generate generates random avatar
func Generate(gender Gender) (image.Image, error) {
	return GenerateForUsername(gender, time.Now().String())
}

// GenerateFile generates random avatar and save it to specified file.
// Image format depends on file extension (jpeg, jpg, png, gif). Default is png
func GenerateFile(gender Gender, filePath string) error {
	img, err := Generate(gender)
	if err != nil {
		return err
	}
	return saveToFile(img, filePath)
}

// GenerateForUsername generates avatar for username
func GenerateForUsername(gender Gender, username string) (image.Image, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(username))
	if err != nil {
		return nil, err
	}
	switch gender {
	case MALE:
		return randomAvatar(assetsStore.Male, int64(h.Sum32()))
	case FEMALE:
		return randomAvatar(assetsStore.Female, int64(h.Sum32()))
	default:
		return nil, ErrUnsupportedGender
	}
}

// GenerateFileForUsername generates avatar for username and save it to specified file.
// Image format depends on file extension (jpeg, jpg, png, gif). Default is png
func GenerateFileForUsername(gender Gender, username string, filePath string) error {
	img, err := GenerateForUsername(gender, username)
	if err != nil {
		return err
	}
	return saveToFile(img, filePath)
}

func saveToFile(img image.Image, filePath string) error {
	outFile, err := os.Create(filePath)
	defer outFile.Close()
	if err != nil {
		return err
	}

	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
	case ".gif":
		err = gif.Encode(outFile, img, nil)
	default:
		err = png.Encode(outFile, img)
	}
	return err
}

func randomAvatar(p person, seed int64) (image.Image, error) {
	rnd := rand.New(rand.NewSource(seed))
	avatar := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var err error
	err = drawImg(avatar, randStringSliceItem(rnd, assetsStore.Background), err)
	err = drawImg(avatar, randStringSliceItem(rnd, p.Face), err)
	err = drawImg(avatar, randStringSliceItem(rnd, p.Clothes), err)
	err = drawImg(avatar, randStringSliceItem(rnd, p.Mouth), err)
	err = drawImg(avatar, randStringSliceItem(rnd, p.Hair), err)
	err = drawImg(avatar, randStringSliceItem(rnd, p.Eye), err)
	return avatar, err
}

func drawImg(dst draw.Image, asset string, err error) error {
	if err != nil {
		return err
	}
	src, _, err := image.Decode(bytes.NewReader(mustAsset(asset)))
	if err != nil {
		return err
	}
	draw.Draw(dst, dst.Bounds(), src, image.Point{0, 0}, draw.Over)
	return nil
}

func mustAsset(fPath string) []byte {
	b, err := fs.ReadFile(dataFS, fPath)
	if err != nil {
		panic(err)
	}
	return b
}

func getPerson(gender Gender) person {
	var genderPath string

	switch gender {
	case FEMALE:
		genderPath = "female"
	case MALE:
		genderPath = "male"
	}

	return person{
		Clothes: mustAssetsList("data/" + genderPath + "/clothes"),
		Eye:     mustAssetsList("data/" + genderPath + "/eye"),
		Face:    mustAssetsList("data/" + genderPath + "/face"),
		Hair:    mustAssetsList("data/" + genderPath + "/hair"),
		Mouth:   mustAssetsList("data/" + genderPath + "/mouth"),
	}
}

func mustAssetsList(dir string) []string {
	dirEntries, err := fs.ReadDir(dataFS, dir)
	if err != nil {
		panic(err)
	}
	assets := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		assets[i] = filepath.Join(dir, dirEntry.Name())
	}
	sort.Sort(naturalSort(assets))
	return assets
}
