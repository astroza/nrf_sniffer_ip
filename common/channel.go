/* (c) 2018 felipe@astroza.cl - See LICENSE
 */
package common

import "io"

type ReadWriterChannel struct {
	In       chan []byte
	Error    chan error
	inBuffer []byte
	Handle   io.ReadWriter
}

func (r *ReadWriterChannel) Init(handle io.ReadWriter) {
	r.In = make(chan []byte)
	r.inBuffer = make([]byte, 4096)
	r.Handle = handle
	go func(reader io.Reader) {
		for {
			bytes, err := reader.Read(r.inBuffer)
			if err == nil {
				r.In <- r.inBuffer[0:bytes]
			} else {
				r.Error <- err
				break
			}
		}
	}(handle)
}
