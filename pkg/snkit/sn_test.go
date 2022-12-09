package snkit

import (
	"testing"
	"time"
)

func TestEncodeTime(t *testing.T) {
	now := time.Now()
	n := uint64(now.Unix())
	t.Log(`ts    :`, n)

	encoded := EncodeBit(n, 5)
	t.Log("result:", encoded)

	input := encoded
	t.Log(DecodeBit(input, 5))

	t.Log(DecodeBit(`VVVVVVV`, 5))
}

func TestEncode(t *testing.T) {
	now := time.Now()
	year, month, day := now.Date()
	sec := now.Sub(now.Truncate(time.Hour * 24)).Seconds()

	t.Log(`ts              :`, now.Unix())
	t.Log("year limit(2022):", EncodeMod(2022%1000, 36))
	t.Log("year limit(1000):", EncodeMod(1000, 36))
	t.Log("day limit(366)  :", EncodeMod(366, 36))
	t.Logf("md limit        : %s%s\n", EncodeMod(12, 36), EncodeMod(31, 36))
	t.Log("sec limit(86400):", EncodeMod(86400, 36))

	t.Log(`now by sec      :`, EncodeMod(uint64(now.Unix()), 36))
	t.Log(`day by sec      :`, uint64(sec), EncodeMod(uint64(sec), 36))
	t.Log(`full now        :`, DecodeMod(`ZZZZZZZ`, 36))

	t.Logf(
		"result          : %s%02s%04s\n",
		EncodeMod(uint64(month), 36)+EncodeMod(uint64(day), 36),
		EncodeMod(uint64(year%1000), 36),
		EncodeMod(uint64(sec), 36),
	)
}
