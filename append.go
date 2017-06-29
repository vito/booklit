package booklit

func Append(a, b Content) Content {
	if a == nil {
		return b
	}

	switch v := a.(type) {
	case nil:
		return b
	case Sequence:
		return Sequence(append(v, b))
	default:
		return Sequence([]Content{a, b})
	}
}
