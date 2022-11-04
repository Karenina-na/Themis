package image

import (
	"math/rand"
	"time"
)

// Factory
//
//	@Description: 图片工厂
func Factory() {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(3)
	switch randNum {
	case 0:
		PrintImage1()
	case 1:
		PrintImage2()
	case 2:
		PrintImage3()
	}
}
