// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en": &dictionary{index: enIndex, data: enData},
		"es": &dictionary{index: esIndex, data: esData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"An error occurred":                                   6,
	"Captcha validation failed":                           3,
	"Failed to read request":                              0,
	"Failed to save contact":                              4,
	"Failed to send email":                                5,
	"Failed to validate captcha, please try again later.": 2,
	"Thanks for contacting us":                            8,
	"The request did not pass validation":                 1,
	"[%s] - New contact":                                  7,
}

var enIndex = []uint32{ // 10 elements
	0x00000000, 0x00000017, 0x0000003b, 0x0000006f,
	0x00000089, 0x000000a0, 0x000000b5, 0x000000c7,
	0x000000dd, 0x000000f6,
} // Size: 64 bytes

const enData string = "" + // Size: 246 bytes
	"\x02Failed to read request\x02The request did not pass validation\x02Fai" +
	"led to validate captcha, please try again later.\x02Captcha validation f" +
	"ailed\x02Failed to save contact\x02Failed to send email\x02An error occu" +
	"rred\x02[%[1]s] - New contact\x02Thanks for contacting us"

var esIndex = []uint32{ // 10 elements
	0x00000000, 0x0000001b, 0x00000040, 0x0000007a,
	0x0000009f, 0x000000bb, 0x000000d5, 0x000000ea,
	0x00000103, 0x0000011c,
} // Size: 64 bytes

const esData string = "" + // Size: 284 bytes
	"\x02Error al leer la petición\x02La petición no pasó la validación\x02Er" +
	"ror al validar el captcha, por favor intente mas tarde.\x02La validación" +
	" de captcha ha fallado\x02Error al salvar el contacto\x02Error al enviar" +
	" el correo\x02Ha ocurrido un error\x02[%[1]s] - Nuevo contacto\x02Gracia" +
	"s por contactarnos"

	// Total table size 658 bytes (0KiB); checksum: 1E7868C3
