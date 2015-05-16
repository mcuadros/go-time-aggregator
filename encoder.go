package aggregator

import (
	"encoding/json"

	"labix.org/v2/mgo/bson"
)

// GetBSON implements bson.Setter, marshall the a TimeAggregator to a string
// ussing TimeAggregator.Marshal
func (a *TimeAggregator) GetBSON() (interface{}, error) {
	if a == nil || len(a.Values) == 0 {
		return nil, nil
	}

	return a.Marshal(), nil
}

// SetBSON implements bson.Setter, unmarshals a bson.Raw into a TimeAggregator
func (a *TimeAggregator) SetBSON(raw bson.Raw) error {
	if raw.Kind == 10 {
		return bson.SetZero
	}

	var bin []byte
	if err := raw.Unmarshal(&bin); err != nil {
		return err
	}

	return a.Unmarshal(bin)
}

// MarshalJSON implements json.Marshaler, creates a JSON representation of the
// aggregator.
func (a *TimeAggregator) MarshalJSON() ([]byte, error) {
	v := make([][]interface{}, len(a.Values))

	i := 0
	for p, a := range a.Values {
		v[i] = []interface{}{p.ToMap(), a.values}
		i++
	}

	return json.Marshal(v)
}
