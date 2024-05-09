package indexes

import (
	"errors"
	"fmt"
)

func NewOffsetAndSize(offset uint64, size uint64) *OffsetAndSize {
	return &OffsetAndSize{
		Offset: offset,
		Size:   size,
	}
}

// IsValid returns true if the offset and size are valid.
func (oas *OffsetAndSize) IsValid() bool {
	return oas.Offset <= maxUint48 && oas.Size <= maxUint24
}

type OffsetAndSize struct {
	Offset uint64 // uint48, 6 bytes, max 281.5 TB (terabytes)
	Size   uint64 // uint24, 3 bytes, max 16.7 MB (megabytes)
}

// Bytes returns the offset and size as a byte slice.
func (oas OffsetAndSize) Bytes() []byte {
	return append(uint48tob(oas.Offset), uint24tob(uint32(oas.Size))...)
}

// FromBytes parses the offset and size from a byte slice.
func (oas *OffsetAndSize) FromBytes(buf []byte) error {
	if len(buf) != IndexValueSize_CidToOffsetAndSize {
		return errors.New("invalid byte slice length")
	}
	_ = buf[IndexValueSize_CidToOffsetAndSize-1] // bounds check hint to compiler
	oas.Offset = btoUint48(buf[:6])
	oas.Size = uint64(btoUint24(buf[6:]))
	return nil
}

func OffsetAndSizeSliceFromBytes(buf []byte) ([]OffsetAndSize, error) {
	if len(buf)%IndexValueSize_CidToOffsetAndSize != 0 {
		return nil, errors.New("invalid byte slice length")
	}
	oass := make([]OffsetAndSize, len(buf)/IndexValueSize_CidToOffsetAndSize)
	for i := 0; i < len(oass); i++ {
		if err := oass[i].FromBytes(buf[i*IndexValueSize_CidToOffsetAndSize : (i+1)*IndexValueSize_CidToOffsetAndSize]); err != nil {
			return nil, fmt.Errorf("failed to parse offset and size at index %d: %w", i, err)
		}
	}
	return oass, nil
}
