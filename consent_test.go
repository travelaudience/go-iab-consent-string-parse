package consent

import (
	"testing"
	"time"
)

func TestNewUserConsent(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Error(err)
	}
	testData := []struct {
		consent                    string
		expectedAllowedVendors     []int
		expectedForbiddenVendors   []int
		expectedAllowedPurposes    []int
		expectedForbiddenPurposes  []int
		expectedConsentCreated     time.Time
		expectedConsentLastUpdated time.Time
		expectedVersion            int
		expectedCmpID              int
		expectedCmpVersion         int
		expectedLanguage           string
	}{
		{
			"BN5lERiOMYEdiAOAWeFRAAYAAaAAptQ",
			[]int{1, 5, 7, 9},
			[]int{0, 10},
			[]int{2, 21},
			[]int{1},
			time.Date(2017, 04, 17, 23, 56, 25, 0, loc),
			time.Date(2018, 04, 17, 23, 56, 25, 0, loc),
			1,
			14,
			22,
			"FR",
		},
		{
			"BN5lERiOMYEdiAKAWXEND1HoSBE6CAFAApAMgBkIDIgM0AgOJxAnQA",
			[]int{225, 515, 5000},
			[]int{0, 1, 3, 3244},
			[]int{4, 24},
			[]int{0, 1, 25},
			time.Date(2017, 04, 17, 23, 56, 25, 0, loc),
			time.Date(2018, 04, 17, 23, 56, 25, 0, loc),
			1,
			10,
			22,
			"EN",
		},
		{
			"BON2bQyON2bQyABABAPLATAAAAAAAA",
			[]int{},
			[]int{0, 1, 2, 3},
			[]int{},
			[]int{0, 1, 2, 3},
			time.Date(2018, 05, 16, 15, 1, 18, 0, loc),
			time.Date(2018, 05, 16, 15, 1, 18, 0, loc),
			1,
			1,
			1,
			"PL",
		},
	}

	for idx, tt := range testData {
		uc, err := NewUserConsent(tt.consent)
		if err != nil {
			t.Error(err)
			return
		}

		for _, v := range tt.expectedAllowedVendors {
			if !uc.IsVendorAllowed(v) {
				t.Errorf("%d. IsVendorAllowed for vendor %d expected true, got false", idx, v)
			}
		}

		for _, v := range tt.expectedForbiddenVendors {
			if uc.IsVendorAllowed(v) {
				t.Errorf("%d. IsVendorAllowed for vendor %d expected false, got true", idx, v)
			}
		}

		if !uc.ArePurposesAllowed(tt.expectedAllowedPurposes) {
			t.Errorf("%d. ArePurposesAllowed for purpose %#v expected true, got false", idx, tt.expectedAllowedPurposes)
		}
		for _, v := range tt.expectedAllowedPurposes {
			if !uc.IsPurposeAllowed(v) {
				t.Errorf("%d. IsPurposeAllowed for purpose %d expected true, got false", idx, v)
			}
		}

		for _, v := range tt.expectedForbiddenPurposes {
			if uc.IsPurposeAllowed(v) {
				t.Errorf("%d. IsPurposeAllowed for purpose %d expected false, got true", idx, v)
			}
		}

		if !uc.ConsentRecordCreated.Equal(tt.expectedConsentCreated) {
			t.Errorf(
				"%d. ConsentRecordCreated expected %s, got %s",
				idx,
				tt.expectedConsentCreated.Format("2006-01-02T15:04:05Z07:00"),
				uc.ConsentRecordCreated.Format("2006-01-02T15:04:05Z07:00"),
			)
		}

		if !uc.ConsentRecordLastUpdated.Equal(tt.expectedConsentLastUpdated) {
			t.Errorf(
				"%d. ConsentRecordLastUpdated expected %s, got %s",
				idx,
				tt.expectedConsentLastUpdated.Format("2006-01-02T15:04:05Z07:00"),
				uc.ConsentRecordLastUpdated.Format("2006-01-02T15:04:05Z07:00"),
			)
		}

		if uc.ConsentString != tt.consent {
			t.Errorf("%d. ConsentString expected %s, got %s", idx, tt.consent, uc.ConsentString)
		}

		if uc.Version != tt.expectedVersion {
			t.Errorf("%d. Version expected %d, got %d", idx, tt.expectedVersion, uc.Version)
		}

		if uc.CmpID != tt.expectedCmpID {
			t.Errorf("%d. CmpID expected %d, got %d", idx, tt.expectedCmpID, uc.CmpID)
		}

		if uc.CmpVersion != tt.expectedCmpVersion {
			t.Errorf("%d. CmpVersion expected %d, got %d", idx, tt.expectedCmpVersion, uc.CmpVersion)
		}

		if uc.ConsentLanguage != tt.expectedLanguage {
			t.Errorf("%d. ConsentLanguage expected %s, got %s", idx, tt.expectedLanguage, uc.ConsentLanguage)
		}
	}
}

func BenchmarkNewUserConsent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := NewUserConsent("BOO1H6gOPmM_3ABABAENBB-AAAAcR7_______9______9uz_Gv_r_f__3nW8_39P_h_7_O__7m_-zzV48_lrQV1yPA1CiIAAAAAAAAAAAA")
		if err != nil {
			b.Error(err)
			return
		}
	}
}