package log

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"unicode/utf8"

	buf "go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const _hex = "0123456789abcdef"

var (
	_pool = buf.NewPool()
)

var _vsvPool = sync.Pool{New: func() interface{} {
	return &csvEncoder{}
}}

func getCsvEncoder() *csvEncoder {
	return _vsvPool.Get().(*csvEncoder)
}

func putCsvEncoder(enc *csvEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_vsvPool.Put(enc)
}

// csvEncoder is a custom zapcore.Encoder, used to output CSV formatted logs.
type csvEncoder struct {
	*zapcore.EncoderConfig
	buf            *buf.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buf.Buffer
	reflectEnc *json.Encoder
}

// NewCSVEncoder creates a new csvEncoder.
func NewCSVEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &csvEncoder{
		EncoderConfig: &cfg,
		buf:           _pool.Get(),
		spaced:        true,
	}
}

// AddArray adds an array to the log entry. csvEncoder does not support arrays.
func (enc *csvEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	// csvEncoder does not support arrays, so we simply ignore this call.
	return nil
}

// AddObject adds an object to the log entry. csvEncoder does not support objects.
func (enc *csvEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	// csvEncoder does not support objects, so we simply ignore this call.
	return nil
}

// AddBinary adds a binary field to the log entry.
func (enc *csvEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

// AddByteString adds a byte string field to the log entry.
func (enc *csvEncoder) AddByteString(key string, val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

// AddBool adds a bool field to the log entry.
func (enc *csvEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

// AddComplex128 adds a complex128 field to the log entry.
func (enc *csvEncoder) AddComplex128(key string, val complex128) {
	// csvEncoder does not support complex numbers, so we simply ignore this call.
	return
}

// AddDuration adds a duration field to the log entry.
func (enc *csvEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(int64(val))
}

// AddFloat64 adds a float64 field to the log entry.
func (enc *csvEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendFloat(val, 64)
}

// AddInt64 adds an int64 field to the log entry.
func (enc *csvEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

// AddReflected adds a reflected field to the log entry. csvEncoder does not support reflection.
func (enc *csvEncoder) AddReflected(key string, obj interface{}) error {
	// csvEncoder does not support reflection, so we simply ignore this call.
	return nil
}

// AddString adds a string field to the log entry.
func (enc *csvEncoder) AddString(key string, val string) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(val)
	enc.buf.AppendByte('"')
}

// AddTime adds a time field to the log entry.
func (enc *csvEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(val.UnixNano())
}

// AddUint64 adds a uint64 field to the log entry.
func (enc *csvEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

// EncodeEntry encodes the log entry.
func (enc *csvEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buf.Buffer, error) {
	final := enc.clone()

	// Add Level as the first field.
	enc.addKey("Level")
	enc.addElementSeparator()
	enc.buf.AppendString(ent.Level.String())

	// Add Time as the second field.
	enc.AddTime("Time", ent.Time)

	// Add LoggerName as the third field.
	enc.addKey("LoggerName")
	enc.addElementSeparator()
	enc.buf.AppendString(ent.LoggerName)

	// Add Message as the fourth field.
	enc.addKey("Message")
	enc.addElementSeparator()
	enc.buf.AppendString(ent.Message)

	// Add Caller information as the fifth field.
	enc.addKey("Caller")
	enc.addElementSeparator()
	if ent.Caller.Defined {
		enc.buf.AppendString(ent.Caller.String())
	} else {
		enc.buf.AppendString("-") // or any other placeholder for undefined caller
	}

	// Add Stack trace as the sixth field.
	enc.addKey("Stack")
	enc.addElementSeparator()
	enc.buf.AppendString(ent.Stack)

	for _, field := range fields {
		enc.AddField(field)
	}
	enc.buf.AppendByte('\n')
	ret := final.buf
	putCsvEncoder(final)
	return ret, nil
}
func (enc *csvEncoder) truncate() {
	enc.buf.Reset()
}

// AddField adds a field to the log entry.
func (enc *csvEncoder) AddField(field zapcore.Field) {
	enc.AddString(field.Key, field.String)
}

// Close implements the Encoder interface.
func (enc *csvEncoder) Close() error {
	return nil
}

// Sync implements the Encoder interface.
func (enc *csvEncoder) Sync() error {
	return nil
}

// addKey adds a key to the log entry.
func (enc *csvEncoder) addKey(key string) {
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(',')
}

// addElementSeparator adds an element separator to the log entry.
func (enc *csvEncoder) addElementSeparator() {
	if enc.buf.Len() > 0 {
		enc.buf.AppendByte(',')
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
func (enc *csvEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *csvEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *csvEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *csvEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (enc *csvEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}
func (enc *csvEncoder) clone() *csvEncoder {
	clone := getCsvEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = _pool.Get()
	return clone
}

// AddComplex64 adds a complex64 field to the log entry.
func (enc *csvEncoder) AddComplex64(key string, value complex64) {
	enc.addKey(key)
	enc.addElementSeparator()
	real2 := fmt.Sprintf("%v", real(value))
	imag2 := fmt.Sprintf("%v", imag(value))
	// Format the complex number as "real+imag*i" or "real-imag*i"
	enc.buf.AppendString(fmt.Sprintf("%s%+si", real2, imag2))
}

// AddFloat32 adds a float32 field to the log entry.
func (enc *csvEncoder) AddFloat32(key string, value float32) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendFloat(float64(value), 32)
}

// AddInt adds an int field to the log entry.
func (enc *csvEncoder) AddInt(key string, value int) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(int64(value))
}

// AddInt32 adds an int32 field to the log entry.
func (enc *csvEncoder) AddInt32(key string, value int32) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(int64(value))
}

// AddInt16 adds an int16 field to the log entry.
func (enc *csvEncoder) AddInt16(key string, value int16) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(int64(value))
}

// AddInt8 adds an int8 field to the log entry.
func (enc *csvEncoder) AddInt8(key string, value int8) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendInt(int64(value))
}

// AddUint adds a uint field to the log entry.
func (enc *csvEncoder) AddUint(key string, value uint) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendUint(uint64(value))
}

// AddUint32 adds a uint32 field to the log entry.
func (enc *csvEncoder) AddUint32(key string, value uint32) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendUint(uint64(value))
}

// AddUint16 adds a uint16 field to the log entry.
func (enc *csvEncoder) AddUint16(key string, value uint16) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendUint(uint64(value))
}

// AddUint8 adds a uint8 field to the log entry.
func (enc *csvEncoder) AddUint8(key string, value uint8) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendUint(uint64(value))
}

// AddUintptr adds a uintptr field to the log entry.
func (enc *csvEncoder) AddUintptr(key string, value uintptr) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendUint(uint64(value))
}

// OpenNamespace is a no-op for csvEncoder since CSV does not support namespaces.
func (enc *csvEncoder) OpenNamespace(key string) {
	// CSV does not support namespaces, so this method does nothing.
}
