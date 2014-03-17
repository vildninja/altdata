AltData
=======

While creating a protocol on top of tcp in go, I felt the core library lacked a decent byte buffer for preparing messages to be sent. Currently this package only contains the Buffer type, but this has already saved me a great deal of lines in my current project.

Inspired by Java's ByteBuffer with Flip(), Rewind() and Clear()

Wrapper for encoding/binary.Read and Write

Write and read from net.Conn - returns output from Conn.Read/Write

	buffer.Clear()
	read, err := buffer.ReadFrom(connection)
	buffer.Flip()


	buffer.Flip()
	read, err := buffer.WriteTo(connection)
	buffer.Clear()