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
type StreamBurstReader struct {
	frames []*dataframe.DataFrame
	stream *os.File
	burstRead bool
	counter int
}

func NewStreamBurstReader(stream *os.File) *StreamBurstReader {
	r := new(StreamBurstReader)
	r.stream = stream
	return r
}


func (r *StreamBurstReader) readBurst() error {
	b, err := ioutil.ReadAll(r.stream)
	if err != nil {
		return err
	}
	burst, err := dataframe.UnmarshalDataBurst(b)
	if err != nil {
		return err
	}
	r.frames = burst.Frames
	return nil
}

func (r StreamBurstReader) NextFrame() (*dataframe.DataFrame, error) {
	if !r.burstRead {
		err := r.readBurst()
		if err != nil {
			return nil, err
		}
	}
	if len(r.frames) <= r.counter {
		return nil, io.EOF
	}
	f := r.frames[r.counter]
	r.counter++
	return f, nil
}

type FileReader struct {
	paths []string
	pathCounter int
	streamReader *StreamBurstReader
}

func NewFileReader(paths []string) *FileReader {
	r := new(FileReader)
	r.paths = paths
	return r
}

func (r *FileReader) nextPath() (string, error) {
	if len(r.paths) <= r.pathCounter {
		return "", io.EOF
	}
	p := r.paths[r.pathCounter]
	r.pathCounter++
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
	r.streamReader = sr
	return nil
}

func (r FileReader) NextFrame() (*dataframe.DataFrame, error) {
	if r.streamReader == nil {
		err := r.nextStream()
		if err != nil {
			return nil, err
		}
	}
	f, err := r.streamReader.NextFrame()
	if err == io.EOF {
		err := r.nextStream()
		if err != nil {
			return nil, err
		}
		f, err = r.streamReader.NextFrame()
	}
	return f, err
}
