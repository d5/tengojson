package tengojson_test

//var mustExist = `func(v) { return v }`
//
//var stringify = `export string`
//
//var mustHaveID = `export func(v) {
//	if !is_map(v) || !v["id"] {
//		return error("missing id")
//	}
//}`
//
//var addTimestamp = `export func(v) {
//}`
//
//func Example() {
//	c, err := tengojson.New().
//		ValidateExpr(".name", `is_string`). // "id" must be string
//		Validate(".")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	_ = c
//}
