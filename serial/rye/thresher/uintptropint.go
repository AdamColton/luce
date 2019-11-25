package thresher

import (
	"unsafe"

	"github.com/adamcolton/luce/serial/rye"
)

type uintPtrOpInt struct{}

func (uintPtrOpInt) size(u uintptr) int {
	i := *(*int)(unsafe.Pointer(u))
	return int(rye.SizeCompactInt64(int64(i)))
}
func (uintPtrOpInt) zero(u uintptr) bool {
	return *(*int)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt) marshal(u uintptr, s *rye.Serializer) {
	i := *(*int)(unsafe.Pointer(u))
	s.CompactInt64(int64(i))
}
func (uintPtrOpInt) unmarshal(u uintptr, d *rye.Deserializer) {
	i := int(d.CompactInt64())
	ptr := (*int)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpInt8 struct{}

func (uintPtrOpInt8) size(u uintptr) int {
	return 1
}
func (uintPtrOpInt8) zero(u uintptr) bool {
	return *(*int8)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt8) marshal(u uintptr, s *rye.Serializer) {
	s.Int8(*(*int8)(unsafe.Pointer(u)))
}
func (uintPtrOpInt8) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*int8)(unsafe.Pointer(u)) = d.Int8()
}

type uintPtrOpInt16C struct{}

func (uintPtrOpInt16C) size(u uintptr) int {
	i := *(*int16)(unsafe.Pointer(u))
	return int(rye.SizeCompactInt64(int64(i)))
}
func (uintPtrOpInt16C) zero(u uintptr) bool {
	return *(*int16)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt16C) marshal(u uintptr, s *rye.Serializer) {
	i := *(*int16)(unsafe.Pointer(u))
	s.CompactInt64(int64(i))
}
func (uintPtrOpInt16C) unmarshal(u uintptr, d *rye.Deserializer) {
	i := int16(d.CompactInt64())
	ptr := (*int16)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpInt16 struct{}

func (uintPtrOpInt16) size(u uintptr) int {
	return 2
}
func (uintPtrOpInt16) zero(u uintptr) bool {
	return *(*int16)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt16) marshal(u uintptr, s *rye.Serializer) {
	s.Int16(*(*int16)(unsafe.Pointer(u)))
}
func (uintPtrOpInt16) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*int16)(unsafe.Pointer(u)) = d.Int16()
}

type uintPtrOpInt32C struct{}

func (uintPtrOpInt32C) size(u uintptr) int {
	i := *(*int32)(unsafe.Pointer(u))
	return int(rye.SizeCompactInt64(int64(i)))
}
func (uintPtrOpInt32C) zero(u uintptr) bool {
	return *(*int32)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt32C) marshal(u uintptr, s *rye.Serializer) {
	i := *(*int32)(unsafe.Pointer(u))
	s.CompactInt64(int64(i))
}
func (uintPtrOpInt32C) unmarshal(u uintptr, d *rye.Deserializer) {
	i := int32(d.CompactInt64())
	ptr := (*int32)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpInt32 struct{}

func (uintPtrOpInt32) size(u uintptr) int {
	return 4
}
func (uintPtrOpInt32) zero(u uintptr) bool {
	return *(*int32)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt32) marshal(u uintptr, s *rye.Serializer) {
	s.Int32(*(*int32)(unsafe.Pointer(u)))
}
func (uintPtrOpInt32) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*int32)(unsafe.Pointer(u)) = d.Int32()
}

type uintPtrOpInt64C struct{}

func (uintPtrOpInt64C) size(u uintptr) int {
	i := *(*int64)(unsafe.Pointer(u))
	return int(rye.SizeCompactInt64(i))
}
func (uintPtrOpInt64C) zero(u uintptr) bool {
	return *(*int64)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt64C) marshal(u uintptr, s *rye.Serializer) {
	i := *(*int64)(unsafe.Pointer(u))
	s.CompactInt64(i)
}
func (uintPtrOpInt64C) unmarshal(u uintptr, d *rye.Deserializer) {
	i := d.CompactInt64()
	ptr := (*int64)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpInt64 struct{}

func (uintPtrOpInt64) size(u uintptr) int {
	return 8
}
func (uintPtrOpInt64) zero(u uintptr) bool {
	return *(*int64)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpInt64) marshal(u uintptr, s *rye.Serializer) {
	s.Int64(*(*int64)(unsafe.Pointer(u)))
}
func (uintPtrOpInt64) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*int64)(unsafe.Pointer(u)) = d.Int64()
}

type uintPtrOpUint struct{}

func (uintPtrOpUint) size(u uintptr) int {
	i := *(*uint)(unsafe.Pointer(u))
	return int(rye.SizeCompactUint64(uint64(i)))
}
func (uintPtrOpUint) zero(u uintptr) bool {
	return *(*uint)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint) marshal(u uintptr, s *rye.Serializer) {
	i := *(*uint)(unsafe.Pointer(u))
	s.CompactUint64(uint64(i))
}
func (uintPtrOpUint) unmarshal(u uintptr, d *rye.Deserializer) {
	i := uint(d.CompactUint64())
	ptr := (*uint)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpUint8 struct{}

func (uintPtrOpUint8) size(u uintptr) int {
	return 1
}
func (uintPtrOpUint8) zero(u uintptr) bool {
	return *(*uint8)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint8) marshal(u uintptr, s *rye.Serializer) {
	s.Uint8(*(*uint8)(unsafe.Pointer(u)))
}
func (uintPtrOpUint8) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*uint8)(unsafe.Pointer(u)) = d.Uint8()
}

type uintPtrOpByte struct{}

func (uintPtrOpByte) size(u uintptr) int {
	return 1
}
func (uintPtrOpByte) zero(u uintptr) bool {
	return *(*byte)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpByte) marshal(u uintptr, s *rye.Serializer) {
	s.Byte(*(*byte)(unsafe.Pointer(u)))
}
func (uintPtrOpByte) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*byte)(unsafe.Pointer(u)) = d.Byte()
}

type uintPtrOpUint16C struct{}

func (uintPtrOpUint16C) size(u uintptr) int {
	i := *(*uint16)(unsafe.Pointer(u))
	return int(rye.SizeCompactUint64(uint64(i)))
}
func (uintPtrOpUint16C) zero(u uintptr) bool {
	return *(*uint16)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint16C) marshal(u uintptr, s *rye.Serializer) {
	i := *(*uint16)(unsafe.Pointer(u))
	s.CompactUint64(uint64(i))
}
func (uintPtrOpUint16C) unmarshal(u uintptr, d *rye.Deserializer) {
	i := uint16(d.CompactUint64())
	ptr := (*uint16)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpUint16 struct{}

func (uintPtrOpUint16) size(u uintptr) int {
	return 2
}
func (uintPtrOpUint16) zero(u uintptr) bool {
	return *(*uint16)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint16) marshal(u uintptr, s *rye.Serializer) {
	s.Uint16(*(*uint16)(unsafe.Pointer(u)))
}
func (uintPtrOpUint16) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*uint16)(unsafe.Pointer(u)) = d.Uint16()
}

type uintPtrOpUint32C struct{}

func (uintPtrOpUint32C) size(u uintptr) int {
	i := *(*uint32)(unsafe.Pointer(u))
	return int(rye.SizeCompactUint64(uint64(i)))
}
func (uintPtrOpUint32C) zero(u uintptr) bool {
	return *(*uint32)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint32C) marshal(u uintptr, s *rye.Serializer) {
	i := *(*uint32)(unsafe.Pointer(u))
	s.CompactUint64(uint64(i))
}
func (uintPtrOpUint32C) unmarshal(u uintptr, d *rye.Deserializer) {
	i := uint32(d.CompactUint64())
	ptr := (*uint32)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpUint32 struct{}

func (uintPtrOpUint32) size(u uintptr) int {
	return 4
}
func (uintPtrOpUint32) zero(u uintptr) bool {
	return *(*uint32)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint32) marshal(u uintptr, s *rye.Serializer) {
	s.Uint32(*(*uint32)(unsafe.Pointer(u)))
}
func (uintPtrOpUint32) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*uint32)(unsafe.Pointer(u)) = d.Uint32()
}

type uintPtrOpUint64C struct{}

func (uintPtrOpUint64C) size(u uintptr) int {
	i := *(*uint64)(unsafe.Pointer(u))
	return int(rye.SizeCompactUint64(i))
}
func (uintPtrOpUint64C) zero(u uintptr) bool {
	return *(*uint64)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint64C) marshal(u uintptr, s *rye.Serializer) {
	i := *(*uint64)(unsafe.Pointer(u))
	s.CompactUint64(i)
}
func (uintPtrOpUint64C) unmarshal(u uintptr, d *rye.Deserializer) {
	i := d.CompactUint64()
	ptr := (*uint64)(unsafe.Pointer(u))
	*ptr = i
}

type uintPtrOpUint64 struct{}

func (uintPtrOpUint64) size(u uintptr) int {
	return 8
}
func (uintPtrOpUint64) zero(u uintptr) bool {
	return *(*uint64)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpUint64) marshal(u uintptr, s *rye.Serializer) {
	s.Uint64(*(*uint64)(unsafe.Pointer(u)))
}
func (uintPtrOpUint64) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*uint64)(unsafe.Pointer(u)) = d.Uint64()
}
