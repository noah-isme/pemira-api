package election

import (

	"testing"
	"time"
)

func mustParsePhaseTime(t *testing.T, value string) *time.Time {
	t.Helper()
	tm, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", value, err)
	}
	return &tm
}

func TestValidatePhaseInputs_AllowsSameDayRanges(t *testing.T) {
	phases := []ElectionPhaseInput{
		{Key: PhaseKeyRegistration, StartAt: mustParsePhaseTime(t, "2025-02-10T00:00:00+07:00"), EndAt: mustParsePhaseTime(t, "2025-02-10T12:00:00+07:00")},
		{Key: PhaseKeyVerification, StartAt: mustParsePhaseTime(t, "2025-02-10T12:30:00+07:00"), EndAt: mustParsePhaseTime(t, "2025-02-10T18:00:00+07:00")},
		{Key: PhaseKeyCampaign, StartAt: mustParsePhaseTime(t, "2025-02-11T08:00:00+07:00"), EndAt: mustParsePhaseTime(t, "2025-02-13T18:00:00+07:00")},
		{Key: PhaseKeyQuietPeriod, StartAt: mustParsePhaseTime(t, "2025-02-14T00:00:00+07:00"), EndAt: mustParsePhaseTime(t, "2025-02-14T18:00:00+07:00")},
		{Key: PhaseKeyVoting, StartAt: mustParsePhaseTime(t, "2025-02-15T08:00:00+07:00"), EndAt: mustParsePhaseTime(t, "2025-02-16T02:00:00+07:00")},
		{Key: PhaseKeyRecap, StartAt: mustParsePhaseTime(t, "2025-02-17T08:00:00+07:00"), EndAt: mustParsePhaseTime(t, "2025-02-18T17:00:00+07:00")},
	}

	if err := validatePhaseInputs(phases); err != nil {
		t.Fatalf("expected same-day ranges to be allowed, got error: %v", err)
	}
}
