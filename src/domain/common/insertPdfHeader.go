package common

import (
	"fmt"
	"os"
	"path"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func InsertPdfHeader(m pdf.Maroto, borderAfterHeader bool) {
	cwd, _ := os.Getwd()
	m.RegisterHeader(func() {
		m.SetBorder(false)
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
		m.Line(8, props.Line{Color: color.NewWhite()})
	})
}

func InsertFixedAssetsRequestPdfHeader(m pdf.Maroto, createdAt string) {
	cwd, _ := os.Getwd()
	m.RegisterHeader(func() {
		m.Row(50, func() {
			// m.Col(2, func() {
			m.FileImage(
				path.Join(cwd, "domain", "assets", "servicios_salud_dgo.jpeg"),
				props.Rect{Percent: 50},
			)
			// })
			m.Row(5, func() {
				m.ColSpace(2)
				m.Col(4, func() {
					m.Text("SERVICIOS DE SALUD DE DURANGO", props.Text{Style: consts.Bold, Align: consts.Center, Size: 9})
				})

				m.ColSpace(2)
				m.SetBorder(true)
				m.Col(2, func() {
					m.Text("Código", props.Text{Align: consts.Left, Left: 3, Size: 6, Top: 1.5})
				})
				m.Col(2, func() {
					m.Text("...", props.Text{Align: consts.Left, Left: 3, Size: 6, Top: 1.5})
				})
				m.SetBorder(false)
			})
			m.Row(5, func() {
				m.ColSpace(2)
				m.Col(4, func() {
					m.Text("DIRECCION ADMINISTRATIVA", props.Text{Style: consts.Bold, Align: consts.Center, Size: 9})
				})
				m.ColSpace(2)
				m.SetBorder(true)
				m.Col(2, func() {
					m.Text(fmt.Sprintf("Fecha de emisión: %v", createdAt), props.Text{Align: consts.Left, Left: 1.5, Size: 6, Top: 1.5})
				})
				m.Col(2, func() {
					m.Text(fmt.Sprintf("Fecha de emisión: %v", createdAt), props.Text{Align: consts.Left, Left: 1.5, Size: 6, Top: 1.5})
				})
				m.SetBorder(false)
			})
			m.Row(5, func() {
				m.ColSpace(2)
				m.Col(4, func() {
					m.Text("SUBDIRECCION DE CONTABILIDAD Y PRESUPUESTOS", props.Text{Style: consts.Bold, Align: consts.Center, Size: 9})
				})
				m.ColSpace(2)
				m.SetBorder(true)
				m.Col(4, func() {
					m.Text("Elaboró: JEFE DEL DEPARTAMENTO DE CONTROL DE ACTIVOS", props.Text{Align: consts.Left, Left: 1.5, Size: 7, Top: 1.5})
				})
				m.SetBorder(false)
			})
			m.Row(5, func() {
				m.ColSpace(2)
				m.Col(4, func() {
					m.Text("DPTO. DE CONTROL DE ACTIVOS E INVENTARIOS", props.Text{Style: consts.Bold, Align: consts.Center, Size: 9})
				})
				m.ColSpace(2)
				m.SetBorder(true)
				m.Col(4, func() {
					m.Text("Aprobado por: SUBDIRECCION DE CONTABILIDAD Y PRESUPUESTOS", props.Text{Align: consts.Left, Left: 1.5, Size: 7, Top: 1.5})
				})
				m.SetBorder(false)
			})
		})
	})
	m.SetBorder(false)
}
