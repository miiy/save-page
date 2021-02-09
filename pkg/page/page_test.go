package page

import (
	"testing"
)

var docTypes = []string{
	`

<!DOCTYPE html>
<html lang="en" data-color-mode="light">
  <head>
    <meta charset="utf-8">
`,
`
<!doctype html>
<html>
    <head>
`,
}

func TestPage_htmlPageDoc(t *testing.T) {
	for i, v := range docTypes {
		docType := htmlDocType(v)
		t.Log(docType)
		if docType == "" {
			t.Error(i)
		}
	}

}