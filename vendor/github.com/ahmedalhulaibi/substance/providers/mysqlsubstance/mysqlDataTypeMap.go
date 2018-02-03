package mysqlsubstance



var regexDataTypePatterns = make(map[string]string)

func init(){
	regexDataTypePatterns["bit.*"] = "int64"
	regexDataTypePatterns["bool.*|tinyint\\(1\\)"] = "bool"
	regexDataTypePatterns["tinyint.*"] = "int8"
	regexDataTypePatterns["unsigned\\stinyint.*"] = "uint8"
	regexDataTypePatterns["smallint.*"] = "int16"
	regexDataTypePatterns["unsigned\\ssmallint.*"] = "uint16"
	regexDataTypePatterns["(mediumint.*|int.*)"] = "int32"
	regexDataTypePatterns["unsigned\\s(mediumint.*|int.*)"] = "uint32"
	regexDataTypePatterns["bigint.*"] = "int64"
	regexDataTypePatterns["unsigned\\sbigint.*"] = "uint64"
	regexDataTypePatterns["(unsigned\\s){0,1}(double.*|float.*|dec.*)"] = "float64"
	regexDataTypePatterns["varchar.*|date.*|time.*|year.*|char.*|.*text.*|enum.*|set.*|.*blob.*|.*binary.*"] = "string"
}