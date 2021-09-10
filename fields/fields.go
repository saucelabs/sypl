package fields

// Fields allows to add structured fields to a message.
type Fields map[string]interface{}

//////
// Helpers.
//////

// Copy keys, and values from `src` into `dst`. If `dst` is `nil`, a new `Fields`
// is initialized. If `src` is nil, nothing happens.
func Copy(src, dst Fields) Fields {
	if src == nil {
		return src
	}

	if dst == nil {
		dst = make(Fields)
	}

	for k, v := range src {
		dst[k] = v
	}

	return dst
}
