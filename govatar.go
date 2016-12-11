package govatar

import (
	"bytes"
	"errors"
	"github.com/skarademir/naturalsort"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type person struct {
	Clothes []string
	Eye     []string
	Face    []string
	Hair    []string
	Mouth   []string
}

type data struct {
	Background []string
	Male       person
	Female     person
}

// Data represents assets database
var Data *data

func init() {
	male := getPerson("male")
	female := getPerson("female")
	Data = &data{Background: readAssetsFrom("data/background"), Male: male, Female: female}
	rand.Seed(time.Now().UTC().UnixNano())
}

type Gender int

const (
	MALE Gender = iota
	FEMALE
)

// Generate generates random avatar
func Generate(gender Gender) (image.Image, error) {
	switch gender {
	case FEMALE:
		return randomAvatar(Data.Female)
	default:
		return randomAvatar(Data.Male)
	}
}

// GenerateFile generates random avatar and save it to specified file
func GenerateFile(gender Gender, file string) error {
	img, err := Generate(gender)
	if err != nil {
		return err
	}
	outFile, err := os.Create(file)
	defer outFile.Close()
	if err != nil {
		return err
	}
	switch strings.ToLower(filepath.Ext(file)) {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{80})
	case ".gif":
		err = gif.Encode(outFile, img, nil)
	default:
		err = png.Encode(outFile, img)
	}
	return err
}

// GenerateFromAssets generates avatar from given assets
func GenerateFromAssets(gender Gender, assets []int) (image.Image, error) {
	switch gender {
	case FEMALE:
		return specificAvatar(Data.Female, assets)
	default:
		return specificAvatar(Data.Male, assets)
	}
}

func specificAvatar(p person, assets []int) (image.Image, error) {
	if len(assets) < 6 {
		return nil, errors.New("Wrong assets size")
	}
	if isOut(Data.Background, assets[0]) ||
		isOut(p.Face, assets[1]) ||
		isOut(p.Clothes, assets[2]) ||
		isOut(p.Mouth, assets[3]) ||
		isOut(p.Hair, assets[4]) ||
		isOut(p.Eye, assets[5]) {
		return nil, errors.New("Wrong assets params")
	}
	avatar := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var err error
	err = drawImg(avatar, Data.Background[assets[0]], err)
	err = drawImg(avatar, p.Face[assets[1]], err)
	err = drawImg(avatar, p.Clothes[assets[2]], err)
	err = drawImg(avatar, p.Mouth[assets[3]], err)
	err = drawImg(avatar, p.Hair[assets[4]], err)
	err = drawImg(avatar, p.Eye[assets[5]], err)
	return avatar, err
}

func isOut(slice []string, i int) bool {
	if i < 0 || i > len(slice)-1 {
		return true
	}
	return false
}

func randomAvatar(p person) (image.Image, error) {
	avatar := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var err error
	err = drawImg(avatar, RandString(Data.Background), err)
	err = drawImg(avatar, RandString(p.Face), err)
	err = drawImg(avatar, RandString(p.Clothes), err)
	err = drawImg(avatar, RandString(p.Mouth), err)
	err = drawImg(avatar, RandString(p.Hair), err)
	err = drawImg(avatar, RandString(p.Eye), err)
	return avatar, err
}

func drawImg(dst draw.Image, asset string, err error) error {
	if err != nil {
		return err
	}
	src, _, err := image.Decode(bytes.NewReader(MustAsset(asset)))
	if err != nil {
		return err
	}
	draw.Draw(dst, dst.Bounds(), src, image.Point{0, 0}, draw.Over)
	return nil
}

func getPerson(gender string) person {
	return person{
		Clothes: readAssetsFrom("data/" + gender + "/clothes"),
		Eye:     readAssetsFrom("data/" + gender + "/eye"),
		Face:    readAssetsFrom("data/" + gender + "/face"),
		Hair:    readAssetsFrom("data/" + gender + "/hair"),
		Mouth:   readAssetsFrom("data/" + gender + "/mouth")}
}

func readAssetsFrom(dir string) []string {
	assets, _ := AssetDir(dir)
	for i, asset := range assets {
		assets[i] = filepath.Join(dir, asset)
	}
	sort.Sort(naturalsort.NaturalSort(assets))
	return assets
}
