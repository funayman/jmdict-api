package parser

import (
	"io"
	"os"
	"strings"
	"testing"
)

var kanjiTestData io.Reader
var entryTestData io.Reader

func init() {
	entryBytes := `<!DOCTYPE JMdict [<!ENTITY adj-no "nouns which may take the genitive case particle 'no'"><!ENTITY n "noun (common) (futsuumeishi)">]><JMdict><entry><ent_seq>1171270</ent_seq><k_ele><keb>右翼</keb><ke_pri>ichi1</ke_pri><ke_pri>news1</ke_pri><ke_pri>nf04</ke_pri></k_ele><r_ele><reb>うよく</reb><re_pri>ichi1</re_pri><re_pri>news1</re_pri><re_pri>nf04</re_pri></r_ele><sense><pos>&adj-no;</pos><gloss>right-wing</gloss><gloss xml:lang="fr">aile droite (oiseau, armée, parti politique, base-ball)</gloss><gloss xml:lang="ru">пра́вое крыло́</gloss><gloss xml:lang="ru">пра́вый фланг</gloss><gloss xml:lang="de">die Rechte</gloss><gloss xml:lang="de">rechter Flügel</gloss></sense><sense><pos>&n;</pos><gloss>right field (e.g. in sport)</gloss><gloss>right flank</gloss><gloss>right wing</gloss><gloss xml:lang="de">{Sport}</gloss><gloss xml:lang="de">rechte Flanke</gloss><gloss xml:lang="de">rechter Flügel</gloss></sense></entry></JMDict>`
	entryTestData = strings.NewReader(entryBytes)

	kanjiBytes := `<kanjidic2><!-- Entry for Kanji: 本 --><character><literal>本</literal><codepoint><cp_value cp_type="ucs">672c</cp_value><cp_value cp_type="jis208">43-60</cp_value></codepoint><radical><rad_value rad_type="classical">75</rad_value><rad_value rad_type="nelson_c">2</rad_value></radical><misc><grade>1</grade><stroke_count>5</stroke_count><variant var_type="jis208">52-81</variant><freq>10</freq><jlpt>4</jlpt></misc><dic_number><dic_ref dr_type="nelson_c">96</dic_ref><dic_ref dr_type="nelson_n">2536</dic_ref><dic_ref dr_type="halpern_njecd">3502</dic_ref><dic_ref dr_type="halpern_kkld">2183</dic_ref><dic_ref dr_type="heisig">211</dic_ref><dic_ref dr_type="gakken">15</dic_ref><dic_ref dr_type="oneill_names">212</dic_ref><dic_ref dr_type="oneill_kk">20</dic_ref><dic_ref dr_type="moro" m_vol="6" m_page="0026">14421</dic_ref><dic_ref dr_type="henshall">70</dic_ref><dic_ref dr_type="sh_kk">25</dic_ref><dic_ref dr_type="sakade">45</dic_ref><dic_ref dr_type="jf_cards">61</dic_ref><dic_ref dr_type="henshall3">76</dic_ref><dic_ref dr_type="tutt_cards">47</dic_ref><dic_ref dr_type="crowley">6</dic_ref><dic_ref dr_type="kanji_in_context">37</dic_ref><dic_ref dr_type="busy_people">2.1</dic_ref><dic_ref dr_type="kodansha_compact">1046</dic_ref><dic_ref dr_type="maniette">215</dic_ref></dic_number><query_code><q_code qc_type="skip">4-5-3</q_code><q_code qc_type="sh_desc">0a5.25</q_code><q_code qc_type="four_corner">5023.0</q_code><q_code qc_type="deroo">1855</q_code></query_code><reading_meaning><rmgroup><reading r_type="pinyin">ben3</reading><reading r_type="korean_r">bon</reading><reading r_type="korean_h">본</reading><reading r_type="ja_on">ホン</reading><reading r_type="ja_kun">もと</reading><meaning>book</meaning><meaning>present</meaning><meaning>main</meaning><meaning>true</meaning><meaning>real</meaning><meaning>counter for long cylindrical things</meaning><meaning m_lang="fr">livre</meaning><meaning m_lang="fr">présent</meaning><meaning m_lang="fr">essentiel</meaning><meaning m_lang="fr">origine</meaning><meaning m_lang="fr">principal</meaning><meaning m_lang="fr">réalité</meaning><meaning m_lang="fr">vérité</meaning><meaning m_lang="fr">compteur d'objets allongés</meaning><meaning m_lang="es">libro</meaning><meaning m_lang="es">origen</meaning><meaning m_lang="es">base</meaning><meaning m_lang="es">contador de cosas alargadas</meaning><meaning m_lang="pt">livro</meaning><meaning m_lang="pt">presente</meaning><meaning m_lang="pt">real</meaning><meaning m_lang="pt">verdadeiro</meaning><meaning m_lang="pt">principal</meaning><meaning m_lang="pt">sufixo p/ contagem De coisas longas</meaning></rmgroup><nanori>まと</nanori></reading_meaning></character></kanjidic2>`
	kanjiTestData = strings.NewReader(kanjiBytes)
}

func TestLoadKanjiDic2(t *testing.T) {
	kanji, err := LoadKanjiDic2(kanjiTestData)
	if err != nil {
		t.Fatal(err)
	}

	if len(kanji) != 1 {
		t.Errorf("Length of 'kanji' is %d expected 1", len(kanji))
	}
}

func TestLoadKanjiDic2File(t *testing.T) {
	//Open Kanji Test Data
	kanjiFile, err := os.Open("../data/test-kanji.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer kanjiFile.Close()

	//Load Data
	kanji, err := LoadKanjiDic2(kanjiFile)
	if err != nil {
		t.Fatal(err)
	}

	//Make sure we have our 1 test kanji
	if len(kanji) != 1 {
		t.Errorf("Length of 'kanji' is %d expected 1", len(kanji))
	}

	//Perform tests
	k := kanji[0]
	if k.Literal != "本" {
		t.Errorf("Expected Kanji.Literal to be '本' got: %s\n", k.Literal)
	}
}

func TestLoadJMDictFile(t *testing.T) {
	entryFile, err := os.Open("../data/test-entry.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer entryFile.Close()

	words, err := LoadJMDict(entryFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(words) != 1 {
		t.Errorf("Length of 'words' is %d expected 1", len(words))
	}
}

func TestLoadJMDict(t *testing.T) {
	words, err := LoadJMDict(entryTestData)
	if err != nil {
		t.Fatal(err)
	}

	if len(words) != 1 {
		t.Errorf("Length of 'words' is %d expected 1", len(words))
	}
}
