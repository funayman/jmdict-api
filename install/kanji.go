//Package entries contains structs of elements parsed in the XML files
package install

import "encoding/xml"

/************************************************************************************************
* The goal of this is to successfully Marshal and Unmarshal kanjidic2 in the exact same format.
* This struct is used when format=kanjidic2 is specified
* struct design and comments taken straight out of the KanjiDic2 DTD
* Project Info: http://www.edrdg.org/kanjidic/kanjd2index.html
* DTD Info: http://www.edrdg.org/kanjidic/kanjidic2_dtdh.html
************************************************************************************************/

//Kanji data for each <character>  element
//<!ELEMENT character (literal,codepoint, radical, misc, dic_number?, query_code?, reading_meaning?)*>
type Kanji struct {
	XMLName xml.Name `xml:"character" json:"-"`
	Literal string   `xml:"literal" json:"literal"`

	//CodePoint element states the code of the character in the various character set standards
	//<!ELEMENT codepoint (cp_value+)>
	CodePoint []struct {
		//The cp_value contains the codepoint of the character in a particular standard. The standard will be identified in the cp_type attribute
		//<!ELEMENT cp_value (#PCDATA)>
		Value string `xml:",chardata" json:"value"`

		//The cp_type attribute states the coding standard applying to the element.
		//The values assigned so far are:
		//  jis208 - JIS X 0208-1997 - kuten coding (nn-nn)
		//  jis212 - JIS X 0212-1990 - kuten coding (nn-nn)
		//  jis213 - JIS X 0213-2000 - kuten coding (p-nn-nn)
		//  ucs - Unicode 4.0 - hex coding (4 or 5 hexadecimal digits)
		//<!ATTLIST cp_value cp_type CDATA #REQUIRED>
		Type string `xml:"cp_type,attr" json:"type"`
	} `xml:"codepoint>cp_value" json:"codepoints"`

	//<!ELEMENT radical (rad_value+)>
	Radical []struct {
		//The radical number, in the range 1 to 214. The particular classification type is stated in the rad_type attribute.
		//<!ELEMENT rad_value (#PCDATA)>
		Value int64 `xml:",chardata" json:"value"`

		//The rad_type attribute states the type of radical classification.
		//  classical - as recorded in the KangXi Zidian.
		//  nelson_c - as used in the Nelson "Modern Japanese-English
		//  Character Dictionary" (i.e. the Classic, not the New Nelson).
		//  This will only be used where Nelson reclassified the kanji.
		//<!ATTLIST rad_value rad_type CDATA #REQUIRED>
		Type string `xml:"rad_type,attr" json:"type"`
	} `xml:"radical>rad_value" json:"radicals"`

	//<!ELEMENT misc (grade?, stroke_count+, variant*, freq?, rad_name*, jlpt?)>
	Misc struct {
		//The kanji grade level. 1 through 6 indicates a Kyouiku kanji and the grade in which the kanji is taught in Japanese schools.
		//8 indicates it is one of the remaining Jouyou Kanji to be learned in junior high school, and 9 or 10 indicates it is a Jinmeiyou (for use in names) kanji.
		//<!ELEMENT grade (#PCDATA)>
		Grade int64 `xml:"grade,omitempty" json:"grade,omitempty"`

		//The stroke count of the kanji, including the radical. If more than one, the first is considered the accepted count, while subsequent ones
		//are common miscounts. (See Appendix E. of the KANJIDIC documentation for some of the rules applied when counting strokes in some of the radicals
		//<!ELEMENT stroke_count (#PCDATA)>
		StrokeCount []int64 `xml:"stroke_count" json:"strokeCount"`

		//Either a cross-reference code to another kanji, usually regarded as a variant, or an alternative indexing code for the current kanji.
		//The type of variant is given in the var_type attribute.
		//<!ELEMENT variant (#PCDATA)>
		Variant []struct {
			Value string `xml:",chardata" json:"value"`

			//The var_type attribute indicates the type of variant code. The current values are:
			//  jis208 - in JIS X 0208 - kuten coding
			//  jis212 - in JIS X 0212 - kuten coding
			//  jis213 - in JIS X 0213 - kuten coding
			//    (most of the above relate to "shinjitai/kyuujitai" alternative character glyphs)
			//  deroo - De Roo number - numeric
			//  njecd - Halpern NJECD index number - numeric
			//  s_h - The Kanji Dictionary (Spahn & Hadamitzky) - descriptor
			//  nelson_c - "Classic" Nelson - numeric
			//  oneill - Japanese Names (O'Neill) - numeric
			//  ucs - Unicode codepoint- hex
			//<!ATTLIST variant var_type CDATA #REQUIRED>
			Type string `xml:"var_type,attr" json:"type"`
		} `xml:"variant,omitempyt" json:"variant,omitempty"`

		//A frequency-of-use ranking. The 2,500 most-used characters have a ranking; those characters that lack this field are not ranked.
		//The frequency is a number from 1 to 2,500 that expresses the relative frequency of occurrence of a character in modern Japanese.
		//This is based on a survey in newspapers, so it is biassed towards kanji used in newspaper articles.
		//The discrimination between the less frequently used kanji is not strong.
		//<!ELEMENT freq (#PCDATA)>
		Frequency int64 `xml:"freq,omitempty" json:"freq,omitempty"`

		//	When the kanji is itself a radical and has a name, this element contains the name (in hiragana).
		//<!ELEMENT rad_name (#PCDATA)>
		RadicalName []string `xml:"rad_name,omitempty" json:"radicalName,omitempty"`

		//The (former) Japanese Language Proficiency test level for this kanji.
		//Values range from 1 (most advanced) to 4 (most elementary).
		//This field does not appear for kanji that were not required for any JLPT level.
		//Note that the JLPT test levels changed in 2010, with a new 5-level system (N1 to N5) being introduced.
		//No official kanji lists are available for the new levels. The new levels are regarded as being similar to the old levels except that the old level 2 is now divided between N2 and N3.
		//<!ELEMENT jlpt (#PCDATA)>
		JLPT int64 `xml:"jlpt,omitempty" json:"jlpt,omitempty"`
	} `xml:"misc" json:"misc"`

	//This element contains the index numbers and similar unstructured information such as
	//page numbers in a number of published dictionaries, and instructional books on kanji.
	//<!ELEMENT dic_number (dic_ref+)>
	Dictionary *dictionaryRef `xml:"dic_number"`

	//These codes contain information relating to the glyph, and can be used for finding a required kanji.
	//The type of code is defined by the qc_type attribute.
	//<!ELEMENT query_code (q_code+)>
	QueryCodes *queryCodes `xml:"query_code" json:"queries"`

	//The readings for the kanji in several languages, and the meanings, also in several languages.
	//The readings and meanings are grouped to enable the handling of the situation where the meaning is differentiated by reading.
	//<!ELEMENT reading_meaning (rmgroup*, nanori*)>
	RM *rm `xml:"reading_meaning" json:"readingMeaning"`
}

/***************** ****************************************************
* Local Structs
*	Needed due to an issue with empty parent tags when marshaling XML
* Github issue -- [https://github.com/golang/go/issues/7233]
**********************************************************************/

type dictionaryRef struct {
	Refs []ref `xml:"dic_ref" json:"dicRefs"`
}

type ref struct {
	//Each dic_ref contains an index number.
	//The particular dictionary, etc. is defined by the dr_type attribute.
	//<!ELEMENT dic_ref (#PCDATA)>
	Index string `xml:",chardata" json:"index"`

	//The dr_type defines the dictionary or reference book, etc. to which dic_ref element applies.
	//<!ATTLIST dic_ref dr_type CDATA #REQUIRED>
	Type  string `xml:"dr_type,attr" json:"type"`
	MVol  string `xml:"m_vol,attr,omitempty" json:"movl,omitempty"`
	MPage string `xml:"m_page,attr,omitempty" json:"mpage,omitempty"`
}

type queryCodes struct {
	Queries []query `xml:"q_code"`
}

type query struct {
	//The q_code contains the actual query-code value, according to the qc_type attribute.
	//<!ELEMENT q_code (#PCDATA)>
	Code string `xml:",chardata" json:"code"`

	//The qc_type attribute defines the type of query code. The current values are:
	//  skip -  Halpern's SKIP (System  of  Kanji  Indexing  by  Patterns) code.
	//  sh_desc - the descriptor codes for The Kanji Dictionary (Tuttle 1996) by Spahn and Hadamitzky.
	//  four_corner - the "Four Corner" code for the kanji.
	//  deroo - the codes developed by the late Father Joseph De Roo, and published in  his book "2001 Kanji" (Bonjinsha).
	//  misclass - a possible misclassification of the kanji according to one of the code types.
	Type string `xml:"qc_type,attr" json:"type"`

	//The values of this attribute indicate the type if misclassification:
	//  posn - a mistake in the division of the kanji
	//  stroke_count - a mistake in the number of strokes
	//  stroke_and_posn - mistakes in both division and strokes
	//  stroke_diff - ambiguous stroke counts depending on glyph
	Misclass string `xml:"skip_misclass,attr,omitempty" json:"misclass,omitempty"`
}

type rm struct {
	//<!ELEMENT rmgroup (reading*, meaning*)>

	//The reading element contains the reading or pronunciation of the kanji.
	//<!ELEMENT reading (#PCDATA)>
	Reading []struct {
		Value string `xml:",chardata" json:"value"`

		//The r_type attribute defines the type of reading in the reading element.
		//The current values are:
		//  pinyin - the modern PinYin romanization of the Chinese reading of the kanji
		//    The tones are represented by a concluding digit
		//  korean_r - the romanized form of the Korean reading(s) of the kanji.
		//    The readings are in the (Republic of Korea) Ministry of Education style of romanization
		//  korean_h - the Korean reading(s) of the kanji in hangul.
		//  ja_on - the "on" Japanese reading of the kanji, in katakana.
		//  	Another attribute r_status, if present, will indicate with
		//  	a value of "jy" whether the reading is approved for a
		//  	"Jouyou kanji".
		//    A further attribute on_type, if present,  will indicate with
		//    a value of kan, go, tou or kan'you the type of on-reading.
		//  ja_kun - the "kun" Japanese reading of the kanji, usually in hiragana.
		//  	Where relevant the okurigana is also included separated by a "."
		//    Readings associated with prefixes and suffixes are marked with a "-"
		//    A second attribute r_status, if present, will indicate with a value of "jy" whether the reading is approved for a "Jouyou kanji".
		//<!ATTLIST reading r_type CDATA #REQUIRED>
		Type string `xml:"r_type,attr" json:"type"`

		//<!ATTLIST reading r_status CDATA #IMPLIED>
		Status string `xml:"r_status,attr,omitempty" json:"status,omitempty"`

		//<!ATTLIST reading on_type CDATA #IMPLIED>
		OnType string `xml:"on_type,attr,omitempty" json:"onType,omitempty"`
	} `xml:"rmgroup>reading" json:"readings"`

	//The meaning associated with the kanji.
	//<!ELEMENT meaning (#PCDATA)>
	Meaning []struct {
		Value string `xml:",innerxml" json:"value"`
		//Value string `xml:",chardata"`

		//The m_lang attribute defines the target language of the meaning.
		//It will be coded using the two-letter language code from the ISO 639-1 standard.
		//When absent, the value "en" (i.e. English) is implied.
		Lang string `xml:"m_lang,attr,omitempty" json:"lang,omitempty"`
	} `xml:"rmgroup>meaning" json:"meanings"`

	//Japanese readings that are now only associated with names.
	//<!ELEMENT nanori (#PCDATA)>
	Nanori []string `xml:"nanori" json:"nanori"`
}

//TODO  Add Authors and Dictionary Title
