package picMaker

import (
	"bufio"
	"github.com/FishZe/Go-DDQuery/account"
	"github.com/FishZe/Go-DDQuery/api"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"time"
)

var font *truetype.Font

func getBuildTime() string {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return "未知"
	}
	return info.ModTime().Format("2006-01-02 15:04:05")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func initPic(user account.User) (*image.RGBA, int, int) {
	columnSum := len(user.VupAttentions) / 300
	if len(user.VupAttentions)%300 != 0 {
		columnSum++
	}
	height := 530 + 55*len(user.VupAttentions)
	if columnSum != 1 {
		height = 530 + 55*300
	}
	img := image.NewRGBA(image.Rect(0, 0, 1000*columnSum, height))
	for x := 0; x < 1000*columnSum; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, image.White)
		}
	}
	return img, columnSum, height
}

func pasteFace(user account.User, img *image.RGBA) {
	f, err := os.Open(user.FaceFile)
	if err != nil {
		log.Printf("打开头像文件失败: %v", err)
		return
	}
	face, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	resizeFace := resize.Resize(200, 200, face, resize.Lanczos3)
	draw.Draw(img, image.Rect(50, 50, 250, 250), resizeFace, image.Point{}, draw.Src)
}

func savePic(user account.User, img *image.RGBA) (string, error) {
	route := "data/out/" + strconv.FormatInt(user.UID, 10) + ".png"
	if !pathExists("data/out") {
		err := os.Mkdir("data/out", 0777)
		if err != nil {
			return "", err
		}
	}
	out, err := os.Create(route)
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}(out)
	b := bufio.NewWriter(out)
	err = png.Encode(b, img)
	if err != nil {
		log.Printf("保存图片失败: %v", err)
		return "", err
	}
	err = b.Flush()
	if err != nil {
		return "", err
	}
	return route, nil
}

func loadFont() *truetype.Font {
	if !pathExists("data/font/SourceHanSansSC-VF.ttf") {
		if !pathExists("data/font") {
			err := os.Mkdir("data/font", 0777)
			if err != nil {
				return nil
			}
		}
		api.DownloadFont()
	}
	fontBytes, err := os.ReadFile("./data/font/SourceHanSansSC-VF.ttf")
	if err != nil {
		log.Printf("读取字体文件失败: %v", err)
		return nil
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("解析字体文件失败: %v", err)
		return nil
	}
	return font
}

func writeText(size float64, x, y int, text string, img *image.RGBA) {
	if font == nil {
		font = loadFont()
	}
	f := freetype.NewContext()
	f.SetDPI(100)
	f.SetFont(font)
	f.SetFontSize(size)
	f.SetClip(img.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255}))
	pt := freetype.Pt(x, y)
	_, err := f.DrawString(text, pt)
	if err != nil {
		log.Printf("写入文字失败: %v", err)
		return
	}
}

func pasteRank(x int, y int, rank int, img *image.RGBA) {
	f, err := os.Open("data/icon/l" + strconv.Itoa(rank) + ".png")
	if err != nil {
		log.Printf("Open rank file failed: %v", err)
		return
	}
	rankImg, _, err := image.Decode(f)
	if err != nil {
		log.Printf("Decode rank file failed: %v", err)
		return
	}
	resizeRank := resize.Resize(60, 30, rankImg, resize.Lanczos3)
	draw.Draw(img, image.Rect(x, y, x+60, y+30), resizeRank, image.Point{}, draw.Src)
}

func writeUserInfo(user account.User, img *image.RGBA, height int) {
	writeText(10, 5, 20, "Generate: "+time.Now().Format("2006/01/02 15:04:05")+"  Build: "+getBuildTime(), img)
	pasteRank(270, 80, user.Rank, img)
	writeText(40, 330, 100, user.Name, img)
	writeText(30, 270, 150, "UID: "+strconv.FormatInt(user.UID, 10), img)
	scale := float64(len(user.VupAttentions)) / float64(len(user.Attentions)) * 100
	writeText(50, 270, 230, strconv.FormatFloat(scale, 'f', 2, 64)+"% ("+strconv.Itoa(len(user.VupAttentions))+" / "+strconv.Itoa(len(user.Attentions))+")", img)
	tm := time.Unix(int64(user.RegTime), 0)
	writeText(25, 50, 300, "注册时间: "+tm.Format("2006-01-02 15:04:05"), img)
	writeText(25, 50, 350, "粉丝量: "+strconv.Itoa(user.Fans)+"    硬币数: "+strconv.FormatFloat(user.Coins, 'f', 1, 64)+"    经验: "+strconv.Itoa(user.CurrentExp), img)
	writeText(25, 50, 400, "舰长: "+strconv.Itoa(len(user.UserGuard.Dd[2]))+"    提督: "+strconv.Itoa(len(user.UserGuard.Dd[1]))+"    总督: "+strconv.Itoa(len(user.UserGuard.Dd[0])), img)
	writeText(15, 20, height-30, "Vup数据来源: https://github.com/dd-center", img)
	writeText(15, 20, height-60, "开源地址: https://github.com/FishZe/Go-DDQuery", img)
}

func checkGuard(user *account.User, x int, y int, mid int64, img *image.RGBA) bool {
	for i := 0; i < 3; i++ {
		for _, v := range user.UserGuard.Dd[i] {
			if v == mid {
				f, err := os.Open("data/icon/g" + strconv.Itoa(i) + ".png")
				if err != nil {
					log.Printf("Open file failed: %v", err)
					return false
				}
				face, _, err := image.Decode(f)
				if err != nil {
					panic(err)
				}
				resizeFace := resize.Resize(50, 50, face, resize.Lanczos3)
				draw.Draw(img, image.Rect(x, y, x+50, y+50), resizeFace, image.Point{}, draw.Src)
				return true
			}
		}
	}
	return false
}

func pasteLiving(x int, y int, img *image.RGBA) {
	f, err := os.Open("data/icon/living.png")
	if err != nil {
		log.Printf("Open file failed: %v", err)
		return
	}
	face, _, err := image.Decode(f)
	if err != nil {
		log.Printf("Decode file failed: %v", err)
		return
	}
	resizeFace := resize.Resize(60, 20, face, resize.Lanczos3)
	draw.Draw(img, image.Rect(x, y, x+60, y+20), resizeFace, image.Point{}, draw.Src)
}

func writeAttention(user account.User, colum int, img *image.RGBA) {
	for i := 0; i < colum; i++ {
		for j := 1; j <= 300 && (i*300+j) <= len(user.VupAttentions); j++ {
			v := user.VupAttentions[i*300+j-1]
			txt := v.Uname + " [" + strconv.FormatInt(v.Mid, 10) + " / " + strconv.FormatInt(v.RoomId, 10) + "]"
			x := 55 + i*1000
			y := 430 + (j-1)*55
			isGuard := checkGuard(&user, x, y, v.Mid, img)
			y += 40
			if isGuard {
				x += 50
			}
			if v.Living {
				pasteLiving(x, y-20, img)
				x += 60
			}
			writeText(30, x, y, txt, img)
		}
	}
}

func MkPic(user account.User) (string, error) {
	img, colum, height := initPic(user)
	pasteFace(user, img)
	writeUserInfo(user, img, height)
	writeAttention(user, colum, img)
	route, err := savePic(user, img)
	if err != nil {
		log.Printf("保存图片失败: %v", err)
		return "", err
	}
	return route, nil
}
