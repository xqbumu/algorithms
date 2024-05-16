package main

import (
	"log"
	"os"

	"github.com/nvisal1/go-wav-codec/pkg/decoder"
	"github.com/nvisal1/go-wav-codec/pkg/encoder"
)

const AZURE_WAV_RESOLUTION = 10000000.0

// 通过Azure返回的发音在WAV文件中的偏移量截取对应的音频
func main() {
	// 02.13
	p := "/tmp/input.wav"
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}

	d := decoder.NewDecoder(f)

	err = d.ReadMetadata()
	if err != nil {
		panic(err)
	}

	log.Println(3700000 * float32(d.Metadata.FMT.SampleRate))
	log.Println(3700000.0 * d.Metadata.FMT.SampleRate)

	sampleRate := float32(d.Metadata.FMT.SampleRate)

	lengthTrans := 3700000 / AZURE_WAV_RESOLUTION * sampleRate
	offsetTrans := 15300000 / AZURE_WAV_RESOLUTION * sampleRate

	buf := make([]int, 0)

	d.ReadAudioData(int(offsetTrans), 0)
	b, err := d.ReadAudioData(int(lengthTrans), 1)
	if err != nil {
		panic(err)
	}

	buf = append(buf, b...)

	f.Close()

	f2, err := os.Create("/tmp/output.wav")
	if err != nil {
		panic(err)
	}
	e, err := encoder.NewEncoder(1, d.Metadata.FMT.NumChannels, d.Metadata.FMT.SampleRate, d.Metadata.FMT.BitsPerSample, f2)
	if err != nil {
		panic(err)
	}

	err = e.WriteAudioData(buf, 0)
	if err != nil {
		panic(err)
	}

	err = e.Close()
	if err != nil {
		panic(err)
	}

	f2.Close()

}
