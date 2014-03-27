package main

import (
	"os"
	"encoding/json"

	"github.com/anchor/dataframe"
)

type FrameEncoder interface {
	EncodeFrame(f *dataframe.DataFrame) ([]byte, error)
}

type BurstEncoder interface {
	EncodeBurst(b *dataframe.DataBurst) ([]byte, error)
}

type RawFrameEncoder struct {}

func (e RawFrameEncoder) EncodeFrame(f *dataframe.DataFrame) ([]byte, error) {
	return dataframe.MarshalDataFrame(f)
}

type JsonFrameEncoder struct {}

func (e JsonFrameEncoder) EncodeFrame(f *dataframe.DataFrame) ([]byte, error) {
	return json.Marshal(f)
}

type RawBurstEncoder struct {}

func (e RawBurstEncoder) EncodeBurst(b *dataframe.DataBurst) ([]byte, error) {
	return dataframe.MarshalDataBurst(b)
}

type JsonBurstEncoder struct {}

func (e JsonBurstEncoder) EncodeBurst(b *dataframe.DataBurst) ([]byte, error) {
	return json.Marshal(b)
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
