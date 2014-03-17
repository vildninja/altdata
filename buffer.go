package altdata

import "encoding/binary"
import "io"
import "errors"

// Buffer inspired by the Java ByteBuffer for simpler data serialization.
type Buffer struct {
	data     []byte
	position int
	limit    int
	order    binary.ByteOrder
}

// Creates a new buffer with given capacity. Default byte order is LittleEndian.
func NewBuffer(capacity int) *Buffer {
	buffer := new(Buffer)
	buffer.data = make([]byte, capacity)
	buffer.Clear()
	buffer.order = binary.LittleEndian
	return buffer
}

func (b *Buffer) Capacity() int {
	return cap(b.data)
}

func (b *Buffer) Length() int {
	return b.limit - b.position
}

func (b *Buffer) Resize(newSize int) {
	newSlice := make([]byte, newSize)
	copy(newSlice, b.data)
	b.data = newSlice
	if b.limit > newSize {
		b.limit = newSize
	}
	if b.position > newSize {
		b.position = newSize
	}
}

// Implementing io.Reader.
func (b *Buffer) Read(p []byte) (n int, err error) {
	n = copy(p, b.data[b.position:b.limit])
	b.position += n
	return
}

// Read a single byte if possible
func (b *Buffer) ReadByte() (byte, error) {
	if b.position >= b.limit {
		return 0, errors.New("No more bytes in buffer")
	}
	b.position++
	return b.data[b.position-1], nil
}

// Read as many bytes as possible from reader.
// Returns bytes read, error.
func (b *Buffer) ReadFrom(reader io.Reader) (n int, err error) {
	n, err = reader.Read(b.Bytes())
	b.position += n
	return
}

// Implementing io.Writer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	n = copy(b.data[b.position:b.limit], p)
	b.position += n
	return
}

// Write a single byte to the Buffer is there is space enough
func (b *Buffer) WriteByte(in byte) error {
	if b.position >= b.limit {
		return errors.New("No more space in buffer")
	}
	b.data[b.position] = in
	b.position++
	return nil
}

// Write as many bytes as possible to writer.
// Returns bytes written, error.
func (b *Buffer) WriteTo(writer io.Writer) (n int, err error) {
	n, err = writer.Write(b.Bytes())
	b.position += n
	return
}

// Sets capacity to current position and position to zero.
// Should be called before reading data from the Buffer.
func (b *Buffer) Flip() {
	b.limit = b.position
	b.position = 0
}

// Reset position to zero.
func (b *Buffer) Rewind() {
	b.position = 0
}

// Sets position to zero and limit to capacity.
// Should be called before starting to write data to the Buffer.
func (b *Buffer) Clear() {
	b.limit = cap(b.data)
	b.position = 0
}

// Returns remaining bytes. DOES NOT move the position.
func (b *Buffer) Bytes() []byte {
	return b.data[b.position:b.limit]
}

// Change the position relatively from its current value.
// Will panic if the changed position is < 0 or > limit!
func (b *Buffer) ChangePosition(n int) {
	if b.position+n < 0 || b.position+n > b.limit {
		panic("Buffer position out of range!")
	}
	b.position += n
}

// Set absolute position and limit manually.
// Will panic if position > limit or limit > capacity.
func (b *Buffer) SetManual(position, limit int) {
	if position > limit || limit > cap(b.data) {
		panic("Buffer position or limit out of range!")
	}
	b.limit = limit
	b.position = position
}

// Returns up to n bytes. Panic if n is negative!
func (b *Buffer) ReadBytes(n int) []byte {
	if n < 0 {
		panic("Tryed to read negative number of bytes from Buffer")
	}

	if b.position+n > b.limit {
		n = b.limit - b.position
	}

	data := b.data[b.position : b.position+n]
	b.position += n

	return data
}

// Returns a string of length up to n bytes. Panic if n is negative!
func (b *Buffer) ReadString(n int) string {
	return string(b.ReadBytes(n))
}

// Writes a string to the buffer. Returns n bytes copied.
func (b *Buffer) WriteString(str string) (n int, err error) {
	return b.Write([]byte(str))
}

// Default is LittleEndian.
func (b *Buffer) SetByteOrder(order binary.ByteOrder) {
	b.order = order
}

// Wrapper for encoding/binary.Read.
func (b *Buffer) ReadBinary(data interface{}) error {
	return binary.Read(b, b.order, data)
}

// Wrapper for encoding/binary.Write.
func (b *Buffer) WriteBinary(data interface{}) error {
	return binary.Write(b, b.order, data)
}
