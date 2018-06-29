package consent

import (
	"encoding/base64"
	"time"
)

const (
	vendorEncodingRange = 1

	versionBitOffset        = 0
	versionBitSize          = 6
	createdBitOffset        = 6
	createdBitSize          = 36
	updatedBitOffset        = 42
	updatedBitSize          = 36
	cmpIDOffset             = 78
	cmpIDSize               = 12
	cmpVersionOffset        = 90
	cmpVersionSize          = 12
	consentScreenSizeOffset = 102
	consentScreenSize       = 6
	consentLanguageOffset   = 108
	consentLanguageSize     = 12
	vendorListVersionOffset = 120
	vendorListVersionSize   = 12
	purposesOffest          = 132
	purposesSize            = 24
	maxVendorIDOffset       = 156
	maxVendorIDSize         = 16
	encodingTypeOffset      = 172
	encodingTypeSize        = 1
	vendorBitFieldOffset    = 173
	defaultConsentOffset    = 173
	numEntriesOffset        = 174
	numEntriesSize          = 12
	rangeEntryOffset        = 186
	vendorIDSize            = 16
)

type UserConsent struct {
	ConsentString            string
	Version                  int
	CmpID                    int
	CmpVersion               int
	ConsentScreenID          int
	ConsentRecordCreated     time.Time
	ConsentRecordLastUpdated time.Time
	ConsentLanguage          string
	VendorListVersion        int

	bits               bits
	purposes           []bool
	vendorEncodingType int
	rangeEntries       []RangeEntry
	defaultConsent     bool
	maxVendorSize      int
}

func NewUserConsent(c string) (UserConsent, error) {
	res, err := base64.RawURLEncoding.DecodeString(c)
	if err != nil {
		return UserConsent{}, err
	}

	b := bits{bytes:res}
	uc := UserConsent{
		ConsentString:     c,
		Version:           b.getInt(versionBitOffset, versionBitSize),
		CmpID:             b.getInt(cmpIDOffset, cmpIDSize),
		CmpVersion:        b.getInt(cmpVersionOffset, cmpVersionSize),
		ConsentScreenID:   b.getInt(consentScreenSizeOffset, consentScreenSize),
		ConsentLanguage:   b.getSixBitString(consentLanguageOffset, consentLanguageSize),
		VendorListVersion: b.getInt(vendorListVersionOffset, vendorListVersionSize),

		maxVendorSize:      b.getInt(maxVendorIDOffset, maxVendorIDSize),
		vendorEncodingType: b.getInt(encodingTypeOffset, encodingTypeSize),
		bits:               b,
	}
	created := b.getInt(createdBitOffset, createdBitSize)
	uc.ConsentRecordCreated = time.Unix(int64(created)/10, 0)
	updated := b.getInt(updatedBitOffset, updatedBitSize)
	uc.ConsentRecordLastUpdated = time.Unix(int64(updated)/10, 0)

	uc.purposes = make([]bool, 0, purposesSize)
	for i, ii := purposesOffest, purposesOffest+purposesSize; i < ii; i++ {
		uc.purposes = append(uc.purposes, b.getBit(i))
	}

	if uc.vendorEncodingType == vendorEncodingRange {
		uc.defaultConsent = b.getBit(defaultConsentOffset)
		numEntries := b.getInt(numEntriesOffset, numEntriesSize)
		currentOffset := rangeEntryOffset
		uc.rangeEntries = make([]RangeEntry, 0, numEntries)
		for i := 0; i < numEntries; i++ {
			rng := b.getBit(currentOffset)
			currentOffset++
			if rng {
				startVendorId := b.getInt(currentOffset, vendorIDSize)
				currentOffset += vendorIDSize
				endVendorId := b.getInt(currentOffset, vendorIDSize)
				currentOffset += vendorIDSize
				uc.rangeEntries = append(uc.rangeEntries, NewRangeEntryWithRange(startVendorId, endVendorId))
			} else {
				vendorId := b.getInt(currentOffset, vendorIDSize)
				currentOffset += vendorIDSize
				uc.rangeEntries = append(uc.rangeEntries, NewRangeEntry(vendorId))
			}
		}
	}

	return uc, nil
}

// IsPurposeAllowed checks if purpose # is allowed for this bid request
// (purpose numeration starts with 1)
func (uc *UserConsent) IsPurposeAllowed(purposeId int) bool {
	if purposeId < 1 || purposeId > len(uc.purposes) {
		return false
	}
	return uc.purposes[purposeId-1]
}

// ArePurposesAllowed checks if purposes are allowed for this bid request
// (purpose numeration starts with 1)
func (uc *UserConsent) ArePurposesAllowed(purposeIds []int) bool {
	for _, p := range purposeIds {
		if p < 1 || p > len(uc.purposes) {
			return false
		}
		if !uc.purposes[p-1] {
			return false
		}
	}
	return true
}

// IsVendorAllowed tells whether vendor is allowed for this bid request
func (uc *UserConsent) IsVendorAllowed(vendorId int) bool {
	if uc.vendorEncodingType == vendorEncodingRange {
		present := uc.findVendorIdInRange(vendorId)
		return present != uc.defaultConsent
	} else {
		return uc.bits.getBit(vendorBitFieldOffset + vendorId - 1)
	}
}

func (uc *UserConsent) findVendorIdInRange(vendorId int) bool {
	limit := len(uc.rangeEntries)
	if limit == 0 {
		return false
	}
	index := limit / 2
	for index >= 0 && index < limit {
		entry := uc.rangeEntries[index]
		if entry.containsVendorId(vendorId) {
			return true
		}
		if index == 0 || index == limit-1 {
			return false
		}
		if entry.idIsGreaterThanMax(vendorId) {
			index = index + ((limit - index) / 2)
		} else {
			index = index / 2
		}
	}
	return false
}

type bits struct {
	bytes    []byte
}

func (b *bits) getInt(startInclusive, size int) int {
	var val int
	sigMask := 1
	var sigIndex = uint(size) - 1

	for i := 0; i < size; i++ {
		if b.getBit(startInclusive + i) {
			val += sigMask << sigIndex
		}
		sigIndex--
	}
	return val
}
func (b *bits) getBit(index int) bool {
	byteIndex := index / 8
	bitOffset := uint(index % 8)
	if byteIndex >= len(b.bytes) {
		return false
	}
	return (b.bytes[byteIndex] & (0x80 >> bitOffset)) != 0
}

func (b *bits) getSixBitString(startInclusive, size int) string {
	if size%6 != 0 {
		return ""
	}
	charNum := size / 6
	str := make([]rune, charNum)
	for i := 0; i < charNum; i++ {
		charCode := b.getInt(startInclusive+(i*6), 6) + 65
		str[i] = rune(charCode)
	}
	return string(str)

}

type RangeEntry struct {
	vendorIds   []int
	maxVendorId int
	minVendorId int
}

func NewRangeEntry(vendorId int) RangeEntry {
	return RangeEntry{
		vendorIds: []int{vendorId},
		maxVendorId: vendorId,
		minVendorId: vendorId,
	}
}

func NewRangeEntryWithRange(startId, endId int) RangeEntry {
	r := RangeEntry{
		vendorIds: make([]int, 0, endId - startId),
		minVendorId: startId,
		maxVendorId: endId,
	}

	for ; startId <= endId; startId++ {
		r.vendorIds = append(r.vendorIds, startId)
	}
	return r
}

func (r *RangeEntry) containsVendorId(vendorId int) bool {
	for _, v := range r.vendorIds {
		if v == vendorId {
			return true
		}
	}
	return false
}

func (r *RangeEntry) idIsGreaterThanMax(vendorId int) bool {
	return vendorId > r.maxVendorId
}
