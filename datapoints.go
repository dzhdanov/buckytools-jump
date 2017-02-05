package buckytools

import (
	"math"
	"sort"
	"time"
)

import "github.com/deniszh/buckytools/whisper"

// FindValidDataPoints does a backwards walk through time to examine the
// highest resolution data for each archive / time period.  We collect valid
// data points and return them in a *[]TimeSeriesPoint.  The second value
// return is an int containing the total number of points examined.  This
// allows one to calculate the percentage of used and unused points stored
// in the Whisper database.
func FindValidDataPoints(wsp *whisper.Whisper) ([]*whisper.TimeSeriesPoint, int, error) {
	points := make([]*whisper.TimeSeriesPoint, 0)
	count := 0
	retentions := whisper.RetentionsByPrecision{wsp.Retentions()}
	sort.Sort(retentions)

	start := int(time.Now().Unix())
	from := 0
	for _, r := range retentions.Iterator() {
		from = int(time.Now().Unix()) - r.MaxRetention()

		ts, err := wsp.Fetch(from, start)
		if err != nil {
			return make([]*whisper.TimeSeriesPoint, 0, 0), 0, err
		}
		count = count + len(ts.Values())
		for _, v := range ts.Points() {
			if !math.IsNaN(v.Value) {
				points = append(points, v)
			}
		}

		start = from
	}

	return points, count, nil
}
