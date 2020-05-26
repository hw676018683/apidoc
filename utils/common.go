package utils

import (
	"reflect"
	"strings"
)

func CommentByTag(tag reflect.StructTag) string {
	comment := tag.Get(`doc`)
	if comment == `` {
		comment = tag.Get(`c`)
	}
	if comment == `` {
		comment = tag.Get(`comment`)
	}
	if strings.Contains(tag.Get(`binding`), `required`) {
		comment = `【必须】` + comment
	}

	return comment
}
