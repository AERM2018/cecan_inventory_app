package common

import (
	"os"
	"path"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func InsertPdfHeader(m pdf.Maroto) {
	cwd, _ := os.Getwd()
	m.RegisterHeader(func() {
		m.Row(15, func() {
			m.Col(1, func() {
				m.FileImage(
					path.Join(cwd, "domain", "assets", "salud_dgo.png"),
					props.Rect{Percent: 100, Top: 1},
				)
			})
			m.Col(1, func() {
				m.FileImage(
					path.Join(cwd, "domain", "assets", "cecan.png"),
					props.Rect{Percent: 100},
				)
			})
			m.Col(10, func() {
				m.Text("Centro de cancerología del estado de Durango", props.Text{
					Size:  18,
					Top:   3,
					Left:  5,
					Style: consts.Bold,
					Align: consts.Left,
				})
			})
		})
	})

}
