package urbanparkjp

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type DecisionCode int

const (
	UNKNOWN DecisionCode = iota
	DECIDE
	UNDECIDE
)

func (dc DecisionCode) String() string {
	switch dc {
	case UNKNOWN:
		return "未確認"
	case DECIDE:
		return "決定"
	case UNDECIDE:
		return "未決定"
	default:
		return ""
	}
}

type ParkTypeCode int

const (
	BLOCK ParkTypeCode = iota + 1
	NEIGHBOR
	COUNTRY
	GENERAL
	SPORTS
	WIDE
	RECREATION
	STATE
	SPECIAL
	BUFFER_GREEN
	URBAN_GREEN
	ROAD_GREEN
	URBAN_FOREST
	SQUARE
)

func (ptc ParkTypeCode) String() string {
	switch ptc {
	case BLOCK:
		return "街区公園"
	case NEIGHBOR:
		return "近隣公園"
	case COUNTRY:
		return "地区公園"
	case GENERAL:
		return "総合公園"
	case SPORTS:
		return "運動公園"
	case WIDE:
		return "広域公園"
	case RECREATION:
		return "レクリエーション都市"
	case STATE:
		return "国営公園"
	case SPECIAL:
		return "特殊公園"
	case BUFFER_GREEN:
		return "緩衝緑地"
	case URBAN_GREEN:
		return "都市緑地"
	case ROAD_GREEN:
		return "緑道"
	case URBAN_FOREST:
		return "都市林"
	case SQUARE:
		return "広場公園"
	default:
		return ""
	}
}

type Dataset struct {
	XMLName xml.Name `xml:"Dataset" json:"xml_name" dynamo:"-"`
	Parks   []Park   `xml:"Park" json:"parks"`
	Points  []Point  `xml:"Point" json:"points"`
}

type Park struct {
	Id      string       `xml:"id,attr" json:"id" dynamo:"id"` // Dinamo pk
	XMLName xml.Name     `xml:"Park" json:"xml_name" dynamo:"-"`
	Loc     Loc          `xml:"loc" json:"loc" dynamo:",set"`
	Adm     string       `xml:"adm" json:"adm"`
	Lgn     string       `xml:"lgn" json:"lgn"`
	Nop     string       `xml:"nop" json:"nop"`
	Kdp     ParkTypeCode `xml:"kdp" json:"kdp"`
	Pop     string       `xml:"pop" json:"pop"`
	Cop     string       `xml:"cop" json:"cop"`
	Opd     Opd          `xml:"opd" json:"opd"`
	Opa     int64        `xml:"opa" json:"opa"`
	Cpd     DecisionCode `xml:"cpd" json:"cpd"`
	Rmk     string       `xml:"rmk" json:"rmk"`
}

type Loc struct {
	Href  string `xml:"href,attr" json:"href" dynamo:"-"`
	Id    string `json:"loc_id"`
	Value Posf64 `json:"loc_value"`
}

type Opd struct {
	TimeInstant TimeInstant `xml:"TimeInstant" json:"time_instant"`
}

type TimeInstant struct {
	TimePosition string `xml:"timePosition" json:"time_position"`
}

type Point struct {
	Id  string `xml:"id,attr" json:"id"`
	Pos string `xml:"pos" json:"pos"`
}

type Posf64 struct {
	Lat float64
	Lon float64
}

func (p *Park) SetLoc(ps map[string]*Posf64) {
	p.Loc.Id = p.Loc.Href[1:]
	pval := ps[p.Loc.Id]
	p.Loc.Value = *pval
}

func SetParksLoc(pks []Park, ps map[string]*Posf64) {
	for i, park := range pks {
		park.SetLoc(ps)
		pks[i] = park
	}
}

func PosToPosf64(pos string) (pf64 *Posf64, err error) {
	fs := strings.Fields(pos)
	lat, err := strconv.ParseFloat(fs[0], 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(fs[1], 64)
	if err != nil {
		return nil, err
	}
	pf64 = &Posf64{Lat: lat, Lon: lon}
	return pf64, nil

}
