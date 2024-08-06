package util

import (
	nanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateID() string {
	return nanoid.MustGenerate("0123456789abcdef", 20)
}

func GenerateMeetingID() string {
	return nanoid.MustGenerate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 8)
}
