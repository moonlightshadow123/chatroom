package utils

import(
	"github.com/astaxie/beego"
	"strings"
)

func Error(errs...interface{}){
	beego.Error(errs)
}

func Info(errs...interface{}){
	beego.Info(errs)
}

func IsImg(f string)bool{
	parts := strings.Split(f, ".")
	suffix := strings.ToLower(parts[len(parts)-1])
	img_suffixes := [4]string{"img", "png", "jpg", "jpeg"}
	for _,img_suf := range img_suffixes{
		if suffix == img_suf{
			return true
		}
	}
	return false
}

func IsVideo(f string)bool{
	parts := strings.Split(f, ".")
	suffix := strings.ToLower(parts[len(parts)-1])
	img_suffixes := [2]string{"mp4", "wav"}
	for _,img_suf := range img_suffixes{
		if suffix == img_suf{
			return true
		}
	}
	return false
}