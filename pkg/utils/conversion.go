package utils

import "encoding/binary"

// Integer constraint for supported integer types
type Integer interface {
	~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64
}

// Float constraint for supported float types
type Float interface {
	~float32 | ~float64
}

// FromBytes converts byte slice to integer (big-endian bytes, big-endian words)
func FromBytes[T Integer](data []byte) T {
	var result T
	switch any(result).(type) {
	case uint16:
		if len(data) >= 2 {
			return T(binary.BigEndian.Uint16(data))
		}
	case int16:
		if len(data) >= 2 {
			return T(int16(binary.BigEndian.Uint16(data)))
		}
	case uint32:
		if len(data) >= 4 {
			return T(binary.BigEndian.Uint32(data))
		}
	case int32:
		if len(data) >= 4 {
			return T(int32(binary.BigEndian.Uint32(data)))
		}
	case uint64:
		if len(data) >= 8 {
			return T(binary.BigEndian.Uint64(data))
		}
	case int64:
		if len(data) >= 8 {
			return T(int64(binary.BigEndian.Uint64(data)))
		}
	}
	return 0
}

// ToBytes converts integer to byte slice (big-endian bytes, big-endian words)
func ToBytes[T Integer](value T) []byte {
	switch any(value).(type) {
	case uint16:
		data := make([]byte, 2)
		binary.BigEndian.PutUint16(data, uint16(value))
		return data
	case int16:
		data := make([]byte, 2)
		binary.BigEndian.PutUint16(data, uint16(value))
		return data
	case uint32:
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data, uint32(value))
		return data
	case int32:
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data, uint32(value))
		return data
	case uint64:
		data := make([]byte, 8)
		binary.BigEndian.PutUint64(data, uint64(value))
		return data
	case int64:
		data := make([]byte, 8)
		binary.BigEndian.PutUint64(data, uint64(value))
		return data
	}
	return nil
}

// Scale applies a scale factor to convert integer to float
func Scale[T Integer, F Float](value T, scale F) F {
	return F(value) * scale
}

// FromBytesWithEndianness converts byte slice to integer using specified byte and word order
func FromBytesWithEndianness[T Integer](data []byte, littleEndianBytes, littleEndianWords bool) T {
	var result T

	switch any(result).(type) {
	case uint16:
		if len(data) >= 2 {
			if littleEndianBytes {
				return T(binary.LittleEndian.Uint16(data))
			}
			return T(binary.BigEndian.Uint16(data))
		}
	case int16:
		if len(data) >= 2 {
			return T(int16(FromBytesWithEndianness[uint16](data, littleEndianBytes, littleEndianWords)))
		}
	case uint32:
		if len(data) >= 4 {
			var reg1, reg2 uint16
			if littleEndianBytes {
				reg1 = binary.LittleEndian.Uint16(data[0:2])
				reg2 = binary.LittleEndian.Uint16(data[2:4])
			} else {
				reg1 = binary.BigEndian.Uint16(data[0:2])
				reg2 = binary.BigEndian.Uint16(data[2:4])
			}

			if littleEndianWords {
				// Low word first, high word second
				return T(uint32(reg1) | (uint32(reg2) << 16))
			} else {
				// High word first, low word second
				return T((uint32(reg1) << 16) | uint32(reg2))
			}
		}
	case int32:
		if len(data) >= 4 {
			return T(int32(FromBytesWithEndianness[uint32](data, littleEndianBytes, littleEndianWords)))
		}
	case uint64:
		if len(data) >= 8 {
			var reg1, reg2, reg3, reg4 uint16
			if littleEndianBytes {
				reg1 = binary.LittleEndian.Uint16(data[0:2])
				reg2 = binary.LittleEndian.Uint16(data[2:4])
				reg3 = binary.LittleEndian.Uint16(data[4:6])
				reg4 = binary.LittleEndian.Uint16(data[6:8])
			} else {
				reg1 = binary.BigEndian.Uint16(data[0:2])
				reg2 = binary.BigEndian.Uint16(data[2:4])
				reg3 = binary.BigEndian.Uint16(data[4:6])
				reg4 = binary.BigEndian.Uint16(data[6:8])
			}

			if littleEndianWords {
				// Words in little-endian order: reg1 (lowest) to reg4 (highest)
				return T(uint64(reg1) | (uint64(reg2) << 16) | (uint64(reg3) << 32) | (uint64(reg4) << 48))
			} else {
				// Words in big-endian order: reg1 (highest) to reg4 (lowest)
				return T((uint64(reg1) << 48) | (uint64(reg2) << 32) | (uint64(reg3) << 16) | uint64(reg4))
			}
		}
	case int64:
		if len(data) >= 8 {
			return T(int64(FromBytesWithEndianness[uint64](data, littleEndianBytes, littleEndianWords)))
		}
	}
	return 0
}

// FromBytesDCBA converts byte slice to integer (little-endian bytes, little-endian words)
func FromBytesDCBA[T Integer](data []byte) T {
	return FromBytesWithEndianness[T](data, true, true)
}

// FromBytesBADC converts byte slice to integer (little-endian bytes, big-endian words)
func FromBytesBADC[T Integer](data []byte) T {
	return FromBytesWithEndianness[T](data, true, false)
}

// FromBytesCDAB converts byte slice to integer (big-endian bytes, little-endian words)
func FromBytesCDAB[T Integer](data []byte) T {
	return FromBytesWithEndianness[T](data, false, true)
}
