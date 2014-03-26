package main

import (
	"os"
	"io"
	"io/ioutil"

	"github.com/anchor/dataframe"
)

type FrameReader interface {
	NextFrame() (*dataframe.DataFrame, error)
}

// StdinReader reads DataFrames from stdin.
// We assume only one burst will be provided on stdin, as bursts
// themselves are the only framing method we use. 
type StreamBurstReaderState struct {
	frames []*dataframe.DataFrame
	burstRead bool
	counter int
}

type StreamBurstReader struct {
	state *StreamBurstReaderState
	stream *os.File
}

func NewStreamBurstReader(stream *os.File) *StreamBurstReader {
	s := new(StreamBurstReaderState)
	r := StreamBurstReader{s, stream}
	return &r
}


func (r StreamBurstReader) readBurst() error {
	b, err := ioutil.ReadAll(r.stream)
	if err != nil {
		return err
	}
	burst, err := dataframe.UnmarshalDataBurst(b)
	if err != nil {
		return err
	}
	r.state.frames = burst.Frames
	return nil
}

func (r StreamBurstReader) NextFrame() (*dataframe.DataFrame, error) {
	if !r.state.burstRead {
		err := r.readBurst()
		r.state.burstRead = true
		if err != nil {
			return nil, err
		}
	}
	if len(r.state.frames) <= r.state.counter {
		return nil, io.EOF
	}
	f := r.state.frames[r.state.counter]
	r.state.counter++
	return f, nil
}

type FileReaderState struct {
	pathCounter int
	streamReader *StreamBurstReader
}

type FileReader struct {
	paths []string
	state *FileReaderState
}

func NewFileReader(paths []string) *FileReader {
	s := new(FileReaderState)
	r := FileReader{paths, s}
	return &r
}

func (r *FileReader) nextPath() (string, error) {
	if len(r.paths) <= r.state.pathCounter {
		return "", io.EOF
	}
	p := r.paths[r.state.pathCounter]
	r.state.pathCounter++
	return p, nil
}

func (r *FileReader) nextStream() error {
	path, err := r.nextPath()
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	sr := NewStreamBurstReader(f)
	r.state.streamReader = sr
	return nil
}

func (r FileReader) NextFrame() (*dataframe.DataFrame, error) {
	if r.state.streamReader == nil {
		err := r.nextStream()
		if err != nil {
			return nil, err
		}
	}
	f, err := r.state.streamReader.NextFrame()
	if err == io.EOF {
		err := r.nextStream()
		if err != nil {
			return nil, err
		}
		f, err = r.state.streamReader.NextFrame()
	}
	return f, err
}
