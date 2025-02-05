package io

import (
	"bufio"
	"io"

	"github.com/juju/ratelimit"
)

type reader struct {
	r      *bufio.Reader
	bucket *ratelimit.Bucket
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Reader(r io.Reader, bucket *ratelimit.Bucket) io.Reader {
	return &reader{
		r:      bufio.NewReader(r),
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}

type writer struct {
	w      *bufio.Writer
	bucket *ratelimit.Bucket
}

// Writer returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Writer(w io.Writer, bucket *ratelimit.Bucket) io.Writer {
	return &writer{
		w:      bufio.NewWriter(w),
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}

type seeker struct {
	s      io.Seeker
	bucket *ratelimit.Bucket
}

// Seeker returns a Seeker that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Seeker(s io.Seeker, bucket *ratelimit.Bucket) io.Seeker {
	return &seeker{
		s:      s,
		bucket: bucket,
	}
}

func (s *seeker) Seek(offset int64, whence int) (int64, error) {
	return s.s.Seek(offset, whence)
}

// // ReadWriter is the interface that groups the basic Read and Write methods.
// type ReadWriter interface {
// 	Reader
// 	Writer
// }

// // ReadCloser is the interface that groups the basic Read and Close methods.
// type ReadCloser interface {
// 	Reader
// 	Closer
// }

// // WriteCloser is the interface that groups the basic Write and Close methods.
// type WriteCloser interface {
// 	Writer
// 	Closer
// }

// // ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
// type ReadWriteCloser interface {
// 	Reader
// 	Writer
// 	Closer
// }

type readseeker struct {
	r reader
	s seeker
}

// ReadSeeker is the interface that groups the basic Read and Seek methods.
func ReadSeeker(r io.Reader, s io.Seeker, bucket *ratelimit.Bucket) io.ReadSeeker {
	return &readseeker{
		r: reader{
			r:      bufio.NewReader(r),
			bucket: bucket,
		},
		s: seeker{
			s:      s,
			bucket: bucket,
		},
	}
}

func (rs *readseeker) Read(buf []byte) (int, error) {
	return rs.r.Read(buf)
}

func (rs *readseeker) Seek(offset int64, whence int) (int64, error) {
	return rs.s.Seek(offset, whence)
}

// // ReadSeekCloser is the interface that groups the basic Read, Seek and Close
// // methods.
// type ReadSeekCloser interface {
// 	Reader
// 	Seeker
// 	Closer
// }

// // WriteSeeker is the interface that groups the basic Write and Seek methods.
// type WriteSeeker interface {
// 	Writer
// 	Seeker
// }

// // ReadWriteSeeker is the interface that groups the basic Read, Write and Seek methods.
// type ReadWriteSeeker interface {
// 	Reader
// 	Writer
// 	Seeker
// }
