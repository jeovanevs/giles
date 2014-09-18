package main

import (
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"strconv"
	"strings"
)

type SmapReading struct {
	Readings [][]interface{}
	UUID     string `json:"uuid"`
}

type SmapMessage struct {
	Readings   *SmapReading
	Metadata   bson.M
	Actuator   bson.M
	Properties bson.M
	UUID       string
	Path       string
}

func (sm *SmapMessage) ToJson() []byte {
	towrite := map[string]*SmapReading{}
	towrite[sm.Path] = sm.Readings
	b, err := json.Marshal(towrite)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return b
}

/*
  We receive the following keys:
  - Metadata: send directly to mongo, if we can
  - Actuator: send directly to mongo, if we can
  - uuid: parse this out for adding to timeseries
  - Readings: parse these out for adding to timeseries
  - Properties: send to mongo, but need to parse out ReadingType to help with parsing Readings
*/
func handleJSON(r io.Reader) ([](*SmapMessage), error) {
	/*
	 * we receive a bunch of top-level keys that we don't know, so we unmarshal them into a
	 * map, and then parse each of the internal objects individually
	 */

	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var ret [](*SmapMessage)
	var e error
	var rawmessage map[string]*json.RawMessage
	err := decoder.Decode(&rawmessage)
	if err != nil {
		return ret, err
	}

	/*
	   global metadata
	   We populate this for every non-endpoint path
	   we come across
	*/
	pathmetadata := make(map[string]interface{})
	isendpoint := true

	for path, reading := range rawmessage {
		isendpoint = true

		js, err := simplejson.NewJson([]byte(*reading))
		if err != nil {
			e = err
		}

		// get uuid
		uuid := js.Get("uuid").MustString("")
		if uuid == "" { // no UUID means no endpoint.
			isendpoint = false
		}

		//get metadata
		localmetadata := js.Get("Metadata").MustMap()

		// if not endpoint, set metadata for this path and then exit
		if !isendpoint {
			pathmetadata[path] = localmetadata
			continue
		}

		message := &SmapMessage{Path: path, UUID: uuid}

		if localmetadata != nil {
			message.Metadata = bson.M(localmetadata)
		}

		// get properties
		properties := js.Get("Properties").MustMap()
		if properties != nil {
			message.Properties = bson.M(properties)
		}

		readingarray := js.Get("Readings").MustArray()
		sr := &SmapReading{UUID: uuid}
		srs := make([][]interface{}, len(readingarray))
		for idx, readings := range readingarray {
			reading := readings.([]interface{})
			ts, _ := strconv.ParseUint(string(reading[0].(json.Number)), 10, 64)
			val, _ := strconv.ParseFloat(string(reading[1].(json.Number)), 64)
			srs[idx] = []interface{}{ts, val}
		}
		sr.Readings = srs
		message.Readings = sr

		actuator := js.Get("Actuator").MustMap()
		if actuator != nil {
			message.Actuator = bson.M(actuator)
		}

		ret = append(ret, message)

	}
	//loop through all path metadata and apply to messages
	for prefix, md := range pathmetadata {
		for idx, msg := range ret {
			if (*msg).Metadata == nil {
				(*msg).Metadata = bson.M(md.(map[string]interface{}))
				ret[idx] = msg
				break
			}
			if strings.HasPrefix((*msg).Path, prefix) {
				for k, v := range md.(map[string]interface{}) {
					if (*msg).Metadata[k] == nil {
						(*msg).Metadata[k] = v
					}
				}
				ret[idx] = msg
			}
		}
	}
	return ret, e
}
