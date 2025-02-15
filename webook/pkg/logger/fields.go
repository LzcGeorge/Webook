package logger

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
