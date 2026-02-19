package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCalculateAge verifies the correctness of calculateAge across boundary conditions
// that the YearDay()-based implementation gets wrong:
//
//   - Leap-year birthday (Feb 29): In non-leap years, Feb 29 does not exist so
//     YearDay() for March 1 (day 60 in a non-leap year) is less than the birthday's
//     YearDay() in a leap year (day 60 as well for Feb 29 in the leap year), causing
//     an off-by-one in certain cross-year calculations.
//   - Birthday today: must NOT subtract 1 (age is already reached today).
//   - Day before birthday: must subtract 1 (birthday has not yet occurred this year).
//   - Day after birthday: must NOT subtract 1.
func TestCalculateAge(t *testing.T) {
	// Use a fixed "today" reference for deterministic tests.
	// We inject via a time.Time parameter by calling calculateAgeAt so tests are
	// not flaky. If calculateAgeAt does not exist yet this test file will not compile,
	// making the suite RED as required.

	tests := []struct {
		name      string
		birthDate time.Time
		now       time.Time
		wantAge   int
	}{
		{
			name:      "standard case - 30 years old",
			birthDate: date(1994, 6, 15),
			now:       date(2024, 6, 15),
			wantAge:   30,
		},
		{
			name:      "birthday today - age equals full years",
			birthDate: date(1990, 3, 1),
			now:       date(2024, 3, 1),
			wantAge:   34,
		},
		{
			name:      "day before birthday - not yet reached this year",
			birthDate: date(1990, 3, 2),
			now:       date(2024, 3, 1),
			wantAge:   33,
		},
		{
			name:      "day after birthday - already passed this year",
			birthDate: date(1990, 2, 28),
			now:       date(2024, 3, 1),
			wantAge:   34,
		},
		// Leap-year birthday: born Feb 29.
		// In non-leap years (e.g. 2023), the birthday threshold is treated as Mar 1
		// for age purposes. The YearDay() approach is wrong because day 60 in a
		// leap year is Feb 29, but day 60 in a non-leap year is Mar 1 — the
		// comparison using YearDay() can therefore produce the wrong result.
		{
			name:    "leap-year birthday (Feb 29) - non-leap check year, before Mar 1",
			// Born Feb 29, 2000. On Feb 28, 2023 they have NOT yet had their
			// birthday this year (birthday treated as Mar 1 in non-leap years).
			birthDate: date(2000, 2, 29),
			now:       date(2023, 2, 28),
			wantAge:   22,
		},
		{
			name:    "leap-year birthday (Feb 29) - non-leap check year, on Mar 1",
			// On Mar 1, 2023 they HAVE had their birthday (Mar 1 threshold).
			birthDate: date(2000, 2, 29),
			now:       date(2023, 3, 1),
			wantAge:   23,
		},
		{
			name:    "leap-year birthday (Feb 29) - in leap year on exact birthday",
			birthDate: date(2000, 2, 29),
			now:       date(2024, 2, 29),
			wantAge:   24,
		},
		{
			name:    "leap-year birthday (Feb 29) - in leap year, one day before",
			birthDate: date(2000, 2, 29),
			now:       date(2024, 2, 28),
			wantAge:   23,
		},
		{
			name:      "just turned 18 today - should pass minimum age check",
			birthDate: date(2006, 1, 15),
			now:       date(2024, 1, 15),
			wantAge:   18,
		},
		{
			name:      "17 years and 364 days - should fail minimum age check",
			birthDate: date(2006, 1, 16),
			now:       date(2024, 1, 15),
			wantAge:   17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAgeAt(tt.birthDate, tt.now)
			assert.Equal(t, tt.wantAge, got,
				"calculateAgeAt(%v, %v) = %d, want %d",
				tt.birthDate.Format("2006-01-02"),
				tt.now.Format("2006-01-02"),
				got,
				tt.wantAge,
			)
		})
	}
}

// date is a helper to construct a time.Time at midnight UTC.
func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
