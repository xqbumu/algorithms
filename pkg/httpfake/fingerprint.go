package httpfake

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type fingerprintParts struct {
	// TODO: debug
	wCsv  *csv.Writer
	count uint

	// serverConn are reused in stream and needs to prevent duplicate parsing.
	printed bool

	settings      map[http2SettingID]uint32
	windowUpdate  uint32
	priorities    []string
	pseudoHeaders []byte
}

func newFingerprintParts() *fingerprintParts {
	fp, err := os.OpenFile(
		path.Join("./tmp", fmt.Sprintf("FPP_%d.csv", time.Now().UnixMicro())),
		os.O_RDWR|os.O_CREATE|os.O_EXCL,
		0600,
	)
	if err != nil {
		panic(err)
	}

	return &fingerprintParts{
		wCsv: csv.NewWriter(fp),
		// the average number of settings here may be 6.
		settings: make(map[http2SettingID]uint32, 6),
		// the average number of priority frame here may be 5.
		priorities: make([]string, 0, 5),
		// any legitimate request will have 3-4 headers.
		pseudoHeaders: make([]byte, 0, 4),
	}
}

// the readFrameResult will no longer exist if readFrames again,
// so it is necessary to save the fingerprint information with plain value.
func (fpp *fingerprintParts) ProcessFrame(res http2readFrameResult) {
	// TODO: debug
	fpp.count++
	fpp.wCsv.Write([]string{
		fmt.Sprintf("`%d", fpp.count),
		fmt.Sprint(res.f.Header().StreamID),
		res.f.Header().Type.String(),
	})

	// once the fingerprint is used, we should not process frame again.
	if fpp.printed {
		return
	}

	// if error occured, the frame will also discard by h2.
	err := res.err
	if err != nil {
		return
	}

	switch f := res.f.(type) {
	case *http2SettingsFrame:
		var sk http2SettingID
		for sk = 0; sk < 6; sk++ {
			if sv, ok := f.Value(sk); ok {
				fpp.settings[sk] = sv
			}
		}
	case *http2WindowUpdateFrame:
		if fpp.windowUpdate > 0 {
			break
		}
		fpp.windowUpdate = f.Increment
	case *http2PriorityFrame:
		fpp.processPriority(f.StreamID, f.http2PriorityParam)
	case *http2MetaHeadersFrame:
		if f.HasPriority() {
			fpp.processPriority(f.StreamID, f.Priority)
		}
		for _, field := range f.Fields {
			if strings.Contains(":method:authority:scheme:path", field.Name) {
				fpp.pseudoHeaders = append(fpp.pseudoHeaders, field.Name[1])
			}
		}
	default:
	}
}

func (fpp *fingerprintParts) processPriority(sid uint32, f http2PriorityParam) {
	fpp.priorities = append(fpp.priorities, fmt.Sprintf("%d:%d:%d:%d", sid, func() uint8 {
		if f.Exclusive {
			return 1
		}
		return 0
	}(), f.StreamDep, f.Weight))
}

func (fpp *fingerprintParts) String() string {
	fpp.printed = true
	fpp.wCsv.Flush()

	buf := bytes.NewBuffer([]byte{})
	var sk http2SettingID
	for sk = 0; sk < 6; sk++ {
		if sv, ok := fpp.settings[sk]; ok {
			fmt.Fprintf(buf, "%d:%d;", sk, sv)
		}
	}
	if len(fpp.settings) > 0 {
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte('|')
	if fpp.windowUpdate == 0 {
		buf.WriteString("00")
	} else {
		fmt.Fprintf(buf, "%d", fpp.windowUpdate)
	}

	buf.WriteByte('|')
	if len(fpp.priorities) == 0 {
		buf.WriteByte('0')
	} else {
		buf.WriteString(strings.Join(fpp.priorities, ","))
	}

	buf.WriteByte('|')
	for k, v := range fpp.pseudoHeaders {
		buf.WriteByte(v)
		if k < len(fpp.pseudoHeaders)-1 {
			buf.WriteByte(',')
		}
	}

	return buf.String()
}
