/* JMDICT */
DROP TABLE IF EXISTS kinf;
DROP TABLE IF EXISTS kpri;
DROP TABLE IF EXISTS audio;
DROP TABLE IF EXISTS rstr;
DROP TABLE IF EXISTS kanj;
DROP TABLE IF EXISTS rinf;
DROP TABLE IF EXISTS rpri;
DROP TABLE IF EXISTS rdng;
DROP TABLE IF EXISTS field;
DROP TABLE IF EXISTS ant;
DROP TABLE IF EXISTS misc;
DROP TABLE IF EXISTS sinf;
DROP TABLE IF EXISTS lsource;
DROP TABLE IF EXISTS gloss;
DROP TABLE IF EXISTS dial;
DROP TABLE IF EXISTS pos;
DROP TABLE IF EXISTS stagk;
DROP TABLE IF EXISTS stagr;
DROP TABLE IF EXISTS xref;
DROP TABLE IF EXISTS sens;
DROP TABLE IF EXISTS enty;

CREATE TABLE enty (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  entseq INTEGER UNIQUE
);

CREATE TABLE kanj (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  kval TEXT,
  eid INTEGER REFERENCES enty (id)
);

/*indicates some unusual aspect, such as okurigana irregularity*/
CREATE TABLE kinf (
  kid INTEGER REFERENCES kanj (id),
  kw TEXT,
  PRIMARY KEY (kid,kw)
);

/*frequency with which the word is used*/
CREATE TABLE kpri (
  kid INTEGER REFERENCES kanj (id),
  kw TEXT,
  PRIMARY KEY (kid,kw)
);

CREATE TABLE rdng (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  rval TEXT,
  eid INTEGER REFERENCES enty (id),
  nokj TEXT
);

CREATE TABLE audio (
  kanj INTEGER NOT NULL REFERENCES kanj (id),
  rdng INTEGER NOT NULL REFERENCES rdng (id),

  /*SHA128 of kanji and reading*/
  name TEXT UNIQUE,
  PRIMARY KEY (kanj,rdng)
);

CREATE TABLE rinf (
  rid INTEGER REFERENCES rdng (id),
  kw TEXT,
  PRIMARY KEY (rid,kw)
);

CREATE TABLE rpri (
  rid INTEGER REFERENCES rdng (id),
  kw TEXT,
  PRIMARY KEY (rid,kw)
);

/*reading restrictions*/
CREATE TABLE rstr (
  kid INTEGER REFERENCES kanj (id),
  rid INTEGER REFERENCES rdng (id),
  PRIMARY KEY (kid,rid)
);

CREATE TABLE sens (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  eid INTEGER REFERENCES enty (id)
);

/*category (e.g. computer, cooking)*/
CREATE TABLE field (
  sid INTEGER REFERENCES sens (id),
  ctg TEXT,
  PRIMARY KEY (sid,ctg)
);

CREATE TABLE ant (
  sid INTEGER REFERENCES sens (id),
  rdng TEXT,
  PRIMARY KEY (sid,rdng)
);

CREATE TABLE misc (
  sid INTEGER REFERENCES sens (id),
  text TEXT,
  PRIMARY KEY (sid,text)
);

CREATE TABLE sinf (
  sid INTEGER REFERENCES sens (id),
  text TEXT,
  PRIMARY KEY (sid,text)
);

CREATE TABLE lsource (
  sid INTEGER REFERENCES sens (id),
  text TEXT,
  lang TEXT DEFAULT 'en',
  type TEXT,
  wasei INTEGER(1) DEFAULT '0',
  PRIMARY KEY (sid,text,lang)
);

CREATE TABLE gloss (
  sid INTEGER REFERENCES sens (id),
  text TEXT,
  lang TEXT,
  gender TEXT,
  PRIMARY KEY (sid,text)
);

/*dialect*/
CREATE TABLE dial (
  sid INTEGER REFERENCES sens (id),
  ben TEXT,
  PRIMARY KEY (sid,ben)
);

CREATE TABLE pos (
  sid INTEGER REFERENCES sens (id),
  kw TEXT,
  PRIMARY KEY (sid,kw)
);

CREATE TABLE stagk (
  sid INTEGER REFERENCES sens (id),
  kval TEXT,
  PRIMARY KEY (sid,kval)
);

CREATE TABLE stagr (
  sid INTEGER REFERENCES sens (id),
  rdng TEXT,
  PRIMARY KEY (sid,rdng)
);

CREATE TABLE xref (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  sid INTEGER REFERENCES sens (id),
  rdng TEXT
);

CREATE INDEX enty_entseq_idx ON enty(entseq);

CREATE INDEX kanj_id_idx ON kanj(id);

CREATE INDEX kanj_kval_idx ON kanj(kval);

CREATE INDEX kanj_eid_idx ON kanj(eid);

CREATE INDEX rdng_id_idx ON rdng(id);

CREATE INDEX rdng_rval_idx ON rdng(rval);

CREATE INDEX rdng_eid_idx ON rdng(eid);

CREATE INDEX sens_id_idx ON sens(id);

CREATE INDEX sens_eid_idx ON sens(eid);

CREATE INDEX field_sid_idx ON field(sid);

CREATE INDEX gloss_sid_idx ON gloss(sid);

CREATE INDEX pos_sid_idx ON pos(sid);

/* KanjiDic2 */
DROP TABLE IF EXISTS codepoint;
DROP TABLE IF EXISTS dictextrainfo;
DROP TABLE IF EXISTS kdictionary;
DROP TABLE IF EXISTS meaning;
DROP TABLE IF EXISTS naori;
DROP TABLE IF EXISTS querycode;
DROP TABLE IF EXISTS radical;
DROP TABLE IF EXISTS radicalName;
DROP TABLE IF EXISTS reading;
DROP TABLE IF EXISTS strokecount;
DROP TABLE IF EXISTS variant;
DROP TABLE IF EXISTS kcharacter;

-- CREATE TABLE kcharacter (
--   id INTEGER PRIMARY KEY AUTOINCREMENT,
--   literal VARCHAR NOT NULL UNIQUE,
--   grade INTEGER,
--   frequency INTEGER,
--   jlpt INTEGER
-- );
--
-- CREATE TABLE codepoint (
--   cid INTEGER REFERENCES kcharacter (id),
--   type VARCHAR,
--   value VARCHAR,
--   PRIMARY KEY (cid,type,value)
-- );
--
-- CREATE TABLE kdictionary (
--   cid INTEGER REFERENCES kcharacter (id),
--   dicindex VARCHAR,
--   name VARCHAR,
--   PRIMARY KEY (cid,dicindex,name)
-- );
--
-- CREATE TABLE dictextrainfo (
--   dict VARCHAR REFERENCES kdictionary (name),
--   volume VARCHAR,
--   page VARCHAR
-- );
--
-- CREATE TABLE meaning (
--   cid INTEGER REFERENCES kcharacter (id),
--   value VARCHAR,
--   lang VARCHAR,
--   insertionOrder INTEGER,
--   PRIMARY KEY (cid,value,lang)
-- );
--
-- CREATE TABLE naori (
--   cid INTEGER REFERENCES kcharacter (id),
--   value VARCHAR,
--   PRIMARY KEY (cid,value)
-- );
--
-- CREATE TABLE querycode (
--   cid INTEGER REFERENCES kcharacter (id),
--   code VARCHAR,
--   type VARCHAR,
--   misclass VARCHAR,
--   PRIMARY KEY (cid,code,type)
-- );
--
-- CREATE TABLE radical (
--   cid INTEGER REFERENCES kcharacter (id),
--   type VARCHAR,
--   value VARCHAR,
--   PRIMARY KEY (cid,type,value)
-- );
--
-- CREATE TABLE radicalName (
--   cid INTEGER REFERENCES kcharacter (id),
--   name VARCHAR
-- );
--
-- CREATE TABLE reading (
--   cid INTEGER REFERENCES kcharacter (id),
--   value VARCHAR,
--   type VARCHAR,
--   status VARCHAR,
--   ontype VARCHAR,
--   PRIMARY KEY (cid,value,type)
-- );
--
-- CREATE TABLE strokecount (
--   cid INTEGER REFERENCES kcharacter (id),
--   count INTEGER,
--   "default" INTEGER
-- );
--
-- CREATE TABLE variant (
--   cid INTEGER REFERENCES kcharacter (id),
--   type VARCHAR,
--   value CHAR,
--   PRIMARY KEY (cid,type,value)
-- );
--
-- CREATE INDEX kcharacter_literal_idx ON kcharacter(literal);


/* VIEWS */
--concat kanjis
DROP VIEW IF EXISTS "vkanjicc";
CREATE VIEW "vkanjicc" AS
SELECT e.id AS "entyid", k.id AS "kanjiid", group_concat(k.kval, "; ") AS "kanjis"
FROM enty e
INNER JOIN kanj k ON k.eid = e.id
GROUP BY e.id;

--concat readings
DROP VIEW IF EXISTS "vreadingcc";
CREATE VIEW "vreadingcc" AS
SELECT e.id AS "entyid", r.id AS "readingid", group_concat(r.rval, "; ") AS "readings"
FROM enty e
INNER JOIN rdng r ON r.eid = e.id
GROUP BY e.id;

--concat glosses
DROP VIEW IF EXISTS "vglosscc";
CREATE VIEW "vglosscc" AS
SELECT s.id AS "senseid", group_concat(g.text, "; ") AS "gloss"
FROM gloss g
INNER JOIN sens s ON s.id = g.sid
GROUP BY s.id;

--concat part-of-speech
DROP VIEW IF EXISTS "vposcc";
CREATE VIEW "vposcc" AS
SELECT p.sid AS "senseid", group_concat(p.kw, "; ") AS "pos"
FROM pos p
LEFT JOIN sens s ON s.id = p.sid
GROUP BY p.sid;

--concat categories (field)
DROP VIEW IF EXISTS "vfieldcc";
CREATE VIEW "vfieldcc" AS
SELECT f.sid AS "senseid", group_concat(f.ctg, "; ") AS "ctg"
FROM field f
INNER JOIN sens s ON f.sid = s.id
GROUP BY s.id;
