package httpfake

import "net/http"

type stdResponseWriter struct {
	rw ResponseWriter
}

func StdResponseWriter(rw ResponseWriter) http.ResponseWriter {
	return &stdResponseWriter{rw}
}

func (std *stdResponseWriter) Header() http.Header {
	return http.Header(std.rw.Header())
}

func (std *stdResponseWriter) Write(b []byte) (int, error) {
	return std.rw.Write(b)
}

func (std *stdResponseWriter) WriteHeader(statusCode int) {
	std.rw.WriteHeader(statusCode)
}

type stdRequest struct {
	r *Request
}

func StdRequest(r *Request) *http.Request {
	ret := &http.Request{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           http.Header(r.Header),
		Body:             r.Body,
		GetBody:          r.GetBody,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          http.Header(r.Trailer),
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
		Cancel:           r.Cancel,
	}

	// ret.Response = StdResponse(ret, r.Response)

	return ret.WithContext(r.Context())
}

type stdResponse struct {
	r *Response
}

func StdResponse(req *http.Request, resp *Response) *http.Response {
	return &http.Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           http.Header(resp.Header),
		Body:             resp.Body,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          http.Header(resp.Trailer),
		Request:          req,
		TLS:              resp.TLS,
	}
}
