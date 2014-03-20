package main

import (
	"github.com/anchor/bletchley/dataframe"
	"os"
)

type FrameEncoder interface {
	EncodeFrame(f *dataframe.DataFrame) ([]byte, error)
}

type BurstEncoder interface {
	EncodeBurst(b *dataframe.DataBurst) ([]byte, error)
}

type RawFrameEncoder struct {
}

func (e RawFrameEncoder) EncodeFrame(f *dataframe.DataFrame) ([]byte, error) {
	return dataframe.MarshalDataFrame(f)
}

type RawBurstEncoder struct {
}

func (e RawBurstEncoder) EncodeFrame(b *dataframe.DataBurst) ([]byte, error) {
	return dataframe.MarshalDataBurst(b)
}

type Writer interface {
	Write(b []byte) error
}

type StreamWriter struct {
	stream *os.File
}

func NewStreamWriter(f *os.File) *StreamWriter {
	s := new(StreamWriter)
	s.stream = f
	return s
}

func (w StreamWriter) Write(b []byte) error {
	_, err := w.stream.Write(b)
	return err
}
