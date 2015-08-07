package aggregator

import (
	"bytes"
	"encoding/gob"
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

func (a *TimeAggregator) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	var err error
	err = enc.Encode(a.kind)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(a.flags)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(uint32(len(a.Values)))
	if err != nil {
		return nil, err
	}

	for p, a := range a.Values {
		err = enc.Encode(p)
		if err != nil {
			return nil, err
		}
		err = enc.Encode(a.values)
		if err != nil {
			return nil, err
		}
		err = enc.Encode(a.p)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (a *TimeAggregator) GobDecode(body []byte) error {
	buf := bytes.NewReader(body)
	dec := gob.NewDecoder(buf)

	var err error
	err = dec.Decode(&a.kind)
	if err != nil {
		return err
	}

	err = dec.Decode(&a.flags)
	if err != nil {
		return err
	}

	a.Values = make(map[Period]*aggregator)

	var length uint32
	err = dec.Decode(&length)
	if err != nil {
		return err
	}

	for i := uint32(0); i < length; i++ {
		var p Period
		err = dec.Decode(&p)
		if err != nil {
			return err
		}

		var agg aggregator
		err = dec.Decode(&agg.values)
		if err != nil {
			return err
		}
		err = dec.Decode(&agg.p)
		if err != nil {
			return err
		}

		a.Values[p] = &agg
	}

	return nil
}
