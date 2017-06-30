package booklit

func Append(first Content, rest ...Content) Content {
	appended := first

	for _, content := range rest {
		if content == nil {
			continue
		}

		switch v := appended.(type) {
		case nil:
			appended = content
		case Sequence:
			return Sequence(append(v, content))
		default:
			return Sequence([]Content{appended, content})
		}
	}

	return appended
}
