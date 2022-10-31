package pdf

import (
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GeneratePDF(files []string, filename string) {
	if len(files) == 0 {
		return
	}
	m := pdf.NewMaroto(consts.Landscape, consts.A4)

	for i := 0; i < len(files); i += 4 {
		fileSlice := files[i:i+4]
		m.Row(80, func() {
			m.Col(6, func() {
				m.FileImage(fileSlice[0], props.Rect{})
			})
			m.Col(6, func() {
				m.FileImage(fileSlice[1], props.Rect{})
			})
		})
		m.Row(3.0, func() {}) // adds spacing between labels
		m.Row(80, func() {
			m.Col(6, func() {
				m.FileImage(fileSlice[2], props.Rect{})
			})
			m.Col(6, func() {
				m.FileImage(fileSlice[3], props.Rect{})
			})
		
		})
		m.AddPage()
	}

	m.OutputFileAndClose(filename)
}
