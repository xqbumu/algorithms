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
	wCsv  *csv.Writer
	count uint

	// serverConn are reused in stream and needs to prevent duplicate parsing.
	printed      bool
	settings     string
	windowUpdate uint32

	// there is no priority frame before headers frame with some browser,
	// such as Microsoft Edge.
	// so it could be empty.
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
		// the average number of priority frame here may be 5.
		priorities: make([]string, 0, 5),
		// any legitimate request will have 3-4 headers.
		pseudoHeaders: make([]byte, 0, 4),
	}
}

// the readFrameResult will no longer exist if readFrames again,
// so it is necessary to save the fingerprint information.
func (fpp *fingerprintParts) ProcessFrame(res http2readFrameResult) {
	// once the fingerprint is used, we should not process frame again.
	if fpp.printed {
		return
	}

	err := res.err
	if err != nil {
		return
	}

	var detail string
	switch f := res.f.(type) {
	case *http2SettingsFrame:
		settings := make([]string, 0, 3)
		for _, k := range []http2SettingID{http2SettingHeaderTableSize, http2SettingInitialWindowSize, http2SettingMaxFrameSize} {
			if v, ok := f.Value(k); ok {
				settings = append(settings, fmt.Sprintf("%d:%d", k, v))
			}
		}
		fpp.settings = strings.Join(settings, ";")
	case *http2WindowUpdateFrame:
		if fpp.windowUpdate > 0 {
			break
		}
		fpp.windowUpdate = f.Increment
	case *http2PriorityFrame:
		detail = fmt.Sprintf("StreamDep: %d", f.StreamDep)
		fpp.priorities = append(fpp.priorities, fmt.Sprintf("%d:%d:%d:%d", f.StreamID, func() uint8 {
			if f.Exclusive {
				return 1
			}
			return 0
		}(), f.StreamDep, f.Weight))
	case *http2MetaHeadersFrame:
		for _, field := range f.Fields {
			if strings.Contains(":method:authority:scheme:path", field.Name) {
				fpp.pseudoHeaders = append(fpp.pseudoHeaders, field.Name[1])
			}
		}
	default:
	}

	fpp.count++
	fpp.wCsv.Write([]string{
		fmt.Sprintf("`%d", fpp.count),
		fmt.Sprint(res.f.Header().StreamID),
		res.f.Header().Type.String(),
		detail,
	})
}

func (fpp *fingerprintParts) String() string {
	fpp.printed = true
	fpp.wCsv.Flush()

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fpp.settings)
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
