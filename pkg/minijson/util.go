package minijson

func MarshalStringMapInferred(s map[string]string) string {
	var jb JsonObjectBuilder
	jb.OpenEx(len(s) * 50)
	for k, v := range s {
		jb.WriteString(k, v)
	}
	jb.Close()
	return jb.String()
}
