package booklit

import (
	"fmt"
	"reflect"
)

type Styled struct {
	Style Style

	Content Content

	Data interface{}
	Flow bool
}

type Style string

const (
	StyleVerbatim    Style = "verbatim"
	StyleItalic      Style = "italic"
	StyleBold        Style = "bold"
	StyleLarger      Style = "larger"
	StyleSmaller     Style = "smaller"
	StyleStrike      Style = "strike"
	StyleSuperscript Style = "superscript"
	StyleSubscript   Style = "subscript"
	StyleInset       Style = "inset"
	StyleAside       Style = "aside"
)

func (con Styled) IsFlow() bool {
	if con.Content != nil {
		return con.Content.IsFlow()
	}

	return con.Flow
}

func (con Styled) String() string {
	return fmt.Sprintf("{styled: %s}", con.Style)
}

func (con Styled) Visit(visitor Visitor) error {
	return visitor.VisitStyled(con)
}

type WalkFunc func(Content) (Content, error)

func (con Styled) Walk(f WalkFunc) (Content, error) {
	if con.Content != nil {
		walked, err := f(con.Content)
		if err != nil {
			return nil, err
		}

		con.Content = walked
	} else {
		walked, err := walkContent(con.Data, f)
		if err != nil {
			return nil, err
		}

		con.Data = walked
	}

	return con, nil
}

func walkContent(data interface{}, f WalkFunc) (interface{}, error) {
	switch v := data.(type) {
	case Content:
		c, err := f(v)
		if err != nil {
			return nil, err
		}

		return c, nil
	default:
		rv := reflect.ValueOf(data)
		rt := rv.Type()

		switch rv.Kind() {
		case reflect.Ptr:
			res, err := walkContent(rv.Elem().Interface(), f)
			if err != nil {
				return nil, err
			}

			return &res, nil
		case reflect.Struct:
			ns := reflect.New(rt).Elem()

			for i := 0; i < rt.NumField(); i++ {
				fv := rv.Field(i)

				walked, err := walkContent(fv.Interface(), f)
				if err != nil {
					return nil, err
				}

				ns.Field(i).Set(reflect.ValueOf(walked))
			}

			return ns.Interface(), nil
		case reflect.Map:
			nm := reflect.MakeMap(rv.Type())

			for _, k := range rv.MapKeys() {
				walked, err := walkContent(rv.MapIndex(k).Interface(), f)
				if err != nil {
					return nil, err
				}

				nm.SetMapIndex(k, reflect.ValueOf(walked))
			}

			return nm.Interface(), nil
		}
	}

	return data, nil
}
