package datetimeconv

import "strings"

const (
	// Slash /
	Slash = "/"
	// Hyphen -
	Hyphen = "-"
)

// DateTimeFormat 日付時間列挙型
type DateTimeFormat string

func datetimePointer(s DateTimeFormat) *DateTimeFormat {
	return &s
}

func (d DateTimeFormat) String() string {
	return string(d)
}

const (
	// YYYY 年
	YYYY = DateTimeFormat("2006")
	// YY 年
	YY = DateTimeFormat("06")
	// MM 月
	MM = DateTimeFormat("01")
	// DD 日
	DD = DateTimeFormat("02")
	// HH 時間
	HH = DateTimeFormat("15")
	// MI 分
	MI = DateTimeFormat("04")
	// SS 秒
	SS = DateTimeFormat("05")
	// HHMISS 150405
	HHMISS = DateTimeFormat(HH + MI + SS)
	// HHMISSColon 15:04:05
	HHMISSColon = DateTimeFormat(HH + ":" + MI + ":" + SS)
	// YYMMDD 060102
	YYMMDD = DateTimeFormat(YY + MM + DD)
	// YYMMDDSlash 06/01/02
	YYMMDDSlash = DateTimeFormat(YY + Slash + MM + Slash + DD)
	// YYMMDDHyphen 06-01-02
	YYMMDDHyphen = DateTimeFormat(YY + Hyphen + MM + Hyphen + DD)
	// YYYYMMDD 20060102
	YYYYMMDD = DateTimeFormat(YYYY + MM + DD)
	// YYYYMMDDSlash 2006/01/02
	YYYYMMDDSlash = DateTimeFormat(YYYY + Slash + MM + Slash + DD)
	// YYYYMMDDHyphen 2006-01-02
	YYYYMMDDHyphen = DateTimeFormat(YYYY + Hyphen + MM + Hyphen + DD)
	// YYYYMMDDhhmiss 20060102150405
	YYYYMMDDhhmiss = DateTimeFormat(YYYYMMDD + HHMISS)
	// YYYYMMDDHhmiss 20060102 150405
	YYYYMMDDHhmiss = DateTimeFormat(YYYYMMDD + " " + HHMISS)
	// YYYYMMDDHhMiSsSlash 2006/01/02 15:04:05
	YYYYMMDDHhMiSsSlash = DateTimeFormat(YYYYMMDDSlash + " " + HHMISSColon)
	// YYYYMMDDHhMiSsHyphen 2006-01-02 15:04:05
	YYYYMMDDHhMiSsHyphen = DateTimeFormat(YYYYMMDDHyphen + " " + HHMISSColon)
	// YYMMDDHhmiss 060102150405
	YYMMDDHhmiss = YY + MM + DD + HH + MI + SS
)

// DateTimeFormatValueOf .
func DateTimeFormatValueOf(s string) *DateTimeFormat {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "yyyy", YYYY.String())
	s = strings.ReplaceAll(s, "yy", YY.String())
	s = strings.ReplaceAll(s, "mm", MM.String())
	s = strings.ReplaceAll(s, "dd", DD.String())
	s = strings.ReplaceAll(s, "hh", HH.String())
	s = strings.ReplaceAll(s, "mi", MI.String())
	s = strings.ReplaceAll(s, "ss", SS.String())
	return datetimePointer(DateTimeFormat(s))
}
