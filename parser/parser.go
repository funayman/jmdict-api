package parser

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/funayman/jmdict-api/entries"
)

//LoadJMDict file and unmarshal all Entries
func LoadJMDict(f io.Reader) (words []*entries.Entry, err error) {
	d, _ := ioutil.ReadAll(f)

	//needed to fix issue
	//https://groups.google.com/forum/#!topic/golang-nuts/yF9RM9rnkYc
	//get all <!ENTITY> objects in XML
	//fix errors when trying to parse &n; &hon; etc
	var rEntity = regexp.MustCompile(`<!ENTITY\s+([^\s]+)\s+"([^"]+)">`)
	entities := make(map[string]string)
	entityDecoder := xml.NewDecoder(bytes.NewReader(d))
	for {
		t, _ := entityDecoder.Token()
		if t == nil {
			break
		}

		dir, ok := t.(xml.Directive)
		if !ok {
			continue
		}

		for _, m := range rEntity.FindAllSubmatch(dir, -1) {
			entities[string(m[1])] = string(m[2])
		}

	}

	decoder := xml.NewDecoder(bytes.NewReader(d)) //go through the data again
	decoder.Entity = entities                     //load entities into the decoder EntityMap
	for {
		//grab all <entry> tokens and Unmarshal into struct
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "entry" {
				var e *entries.Entry

				if err = decoder.DecodeElement(&e, &se); err != nil {
					return nil, err
				}
				words = append(words, e)
			}
		default:
			//do nothing
		}

	}
	return words, nil
}

//LoadKanjiDic2 and return all Kanji
func LoadKanjiDic2(data io.Reader) (characters []*entries.Kanji, err error) {
	xd := xml.NewDecoder(data)
	for {
		t, _ := xd.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "character" {
				var c *entries.Kanji
				if err = xd.DecodeElement(&c, &se); err != nil {
					return
				}
				characters = append(characters, c)
			}
		}
	}

	return
}
