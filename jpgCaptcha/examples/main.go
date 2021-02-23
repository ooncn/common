package main

import (
	"image/color"
	"image/jpeg"
	"net/http"
)

var captcha = jpgCaptcha.New()

func main() {

	_ = captcha.AddFont(util.GetCurrentDirectory() + "/COLONNA.TTF")
	captcha.SetDisturbance(jpgCaptcha.HIGH)
	//设置颜色
	captcha.SetFrontColor(color.Black, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		gifData, code := captcha.RangCaptcha()
		println(code)
		jpeg.Encode(w, gifData, nil)
	})
	http.ListenAndServe(":7080", nil)
}
