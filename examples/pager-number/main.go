package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	// 输入图片文件夹路径和输出图片文件夹路径
	inputFolder := flag.String("input", ".", "输出目录")
	outputFolder := flag.String("output", "output", "输出目录")
	flag.Parse()

	_, err := os.Stat(*outputFolder)
	if err != nil {
		// Directory doesn't exist, create it
		err := os.MkdirAll(*outputFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create directory: %s\n", err)
			return
		}
		fmt.Println("Directory created successfully!")
	} else {
		fmt.Println("Directory already exists!")
	}

	// 获取输入文件夹中的所有图片文件
	imageFiles, err := filepath.Glob(filepath.Join(*inputFolder, "*.jpg"))
	if err != nil {
		fmt.Println("无法读取输入文件夹中的图片文件:", err)
		return
	}

	for _, imagePath := range imageFiles {
		// 读取图片文件
		file, err := os.Open(imagePath)
		if err != nil {
			fmt.Println("无法打开图片文件:", err)
			continue
		}
		defer file.Close()

		// 解码图片
		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Println("无法解码图片:", err)
			continue
		}

		// 创建一个带注释的图像
		annotatedImg := addPageNumber(img, getPageNumber(imagePath))

		// 创建输出文件
		outputPath := filepath.Join(*outputFolder, filepath.Base(imagePath))
		outputFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Println("无法创建输出文件:", err)
			continue
		}
		defer outputFile.Close()

		// 将图像编码为JPEG格式并保存到输出文件
		err = jpeg.Encode(outputFile, annotatedImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Println("无法保存图像:", err)
			continue
		}

		fmt.Println("已处理:", imagePath)
	}

	fmt.Println("所有图片处理完成。")
}

// 添加页码到图像
func addPageNumber(img image.Image, pageNumber int) image.Image {
	// 创建绘图上下文
	dc := gg.NewContextForImage(img)

	// 创建文本标签
	label := fmt.Sprintf("%4d", pageNumber)

	// 计算文本位置
	textWidth, textHeight := dc.MeasureString(label)
	textX := float64(dc.Width()) - textWidth - 30
	textY := 30 + textHeight

	// 设置文本样式
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 48})
	dc.SetFontFace(face)

	// 在图像上绘制文本标签
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(label, textX, textY, 1.0, 0.0)

	return dc.Image()
}

// 获取文件名中的页码
func getPageNumber(filePath string) int {
	filename := filepath.Base(filePath)
	extension := filepath.Ext(filename)
	pageNumberStr := filename[:len(filename)-len(extension)]
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		fmt.Printf("无法解析文件名中的页码：%s\n", filename)
		return 0
	}
	return pageNumber
}
