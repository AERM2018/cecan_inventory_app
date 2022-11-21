package models

import (
	"cecan_inventory/domain/common"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type FixedAssetsRequestDoc struct {
	FixedAssetRequest FixedAssetsRequestDetailed
	CreatedAt         time.Time
	SignaturesInfo    FixedAssetsRequestSignaturesInfo
}

type FixedAssetsRequestSignaturesInfo struct {
	Director      User
	Administrator User
}

func (doc *FixedAssetsRequestDoc) CreateDoc() (string, error) {
	cwd, _ := os.Getwd()
	doc.CreatedAt = time.Now()
	// outputFilePath := path.Join(cwd, "domain", "pdfs", fmt.Sprintf("%v_%v.pdf", "fixedAssetsRequest", doc.CreatedAt.UnixMilli()))
	outputFilePath := path.Join(cwd, "domain", "pdfs", fmt.Sprintf("%v.pdf", "fixedAssetsRequest"))
	mto := pdf.NewMaroto(consts.Landscape, consts.Letter)
	common.InsertFixedAssetsRequestPdfHeader(mto, doc.FixedAssetRequest.CreatedAt.Format(time.RFC3339))
	doc._addFixedAssetRequestInfo(mto)
	doc._addFixedAssetsList(mto)
	doc._addSignatues(mto)
	err := mto.OutputFileAndClose(outputFilePath)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}
	return outputFilePath, nil
}

func (doc FixedAssetsRequestDoc) _addFixedAssetRequestInfo(m pdf.Maroto) {
	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("FORMATO UNICO DE RESGUARDO PARA CONTROL INTERNO DE ACTIVO FIJO EN LAS UNIDADES", props.Text{Size: 11, Style: consts.Bold})
		})
	})
	m.Row(2, func() {
		m.Col(2, func() {
			m.Text("UNIDAD DE ADSCRIPCIÓN:", props.Text{Size: 8, Style: consts.Bold})
		})
		m.Col(3, func() {
			m.Text("CENTRO ESTATAL DE CANCEROLOGÍA", props.Text{Size: 8, Style: consts.Bold, Align: consts.Center})
		})
	})
	m.Row(2, func() {
		m.ColSpace(2)
		m.Col(4, func() {
			m.Text("__________________________________________", props.Text{Size: 8, Style: consts.Bold})
		})
	})
	// --------------------------------
	// This simulates blank space
	m.Row(5, func() {
		m.ColSpace(3)
		m.Col(4, func() {
			m.Text("")
		})
	})
	m.Row(2, func() {
		m.Col(2, func() {
			m.Text("DEPARTAMENTO:", props.Text{Size: 8, Style: consts.Bold})
		})
		m.Col(3, func() {
			m.Text(doc.FixedAssetRequest.Department.Name, props.Text{Size: 8, Style: consts.Bold, Align: consts.Center})
		})
		m.ColSpace(3)
		m.Col(1, func() {
			m.Text("FECHA:", props.Text{Size: 8, Style: consts.Bold})
		})
		m.Col(2, func() {
			fmt.Println(doc.CreatedAt)
			m.Text(doc.CreatedAt.Format("02/01/2006"), props.Text{Size: 8, Style: consts.Bold, Align: consts.Center})
		})
	})
	// This simulates the line below the text
	m.Row(2, func() {
		m.ColSpace(2)
		m.Col(4, func() {
			m.Text("_______________________________________", props.Text{Size: 8, Style: consts.Bold})
		})
		m.ColSpace(3)
		m.Col(1, func() {
			m.Text("________________________", props.Text{Size: 8, Style: consts.Bold})
		})
	})
}

func (doc FixedAssetsRequestDoc) _addFixedAssetsList(m pdf.Maroto) {
	m.Line(5, props.Line{Color: color.NewWhite()})
	gridSizes := []uint{2, 5, 1, 1, 1, 1, 1}
	headers := []string{"ETIQUETA (CLAVE DE INVENTARIO)", "DESCRIPCION", "MARCA", "MODELO", "SERIE", "TIPO", "ESTADO FISICO"}
	contents := make([][]string, 0)
	for _, item := range doc.FixedAssetRequest.FixedAssets {
		contents = append(contents, []string{
			item.FixedAsset.Key,
			item.FixedAsset.Description,
			item.FixedAsset.Brand,
			item.FixedAsset.Model,
			item.FixedAsset.Series,
			item.FixedAsset.Type,
			item.FixedAsset.PhysicState})
	}
	for len(contents) < 10 {
		contents = append(contents, []string{"", "", "", "", "", "", ""})
	}
	m.SetBorder(true)
	m.Row(10, func() {
		for gz, h := range headers {
			m.Col(gridSizes[gz], func() {
				m.Text(h, props.Text{Left: 1.5, Top: 1, Style: consts.Bold, Size: 8})
			})

		}
	})
	for _, row := range contents {
		rowHeight := 7
		if len(row[1]) > 30 {
			rowHeight = 13
		}
		m.Row(float64(rowHeight), func() {
			for gz, col := range row {
				m.Col(gridSizes[gz], func() {
					m.Text(col, props.Text{Left: 1.5, Top: 1, Extrapolate: false, Size: 7})
				})
			}
		})
	}
}

func (doc FixedAssetsRequestDoc) _addSignatues(m pdf.Maroto) {
	m.SetBorder(false)
	m.Line(5, props.Line{Color: color.NewWhite()})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("DE ACUERDO A LOS ARTICULOS ARTICULOS 23 Y 24 DEL CAPITULO II DEL REGISTRO PATRIMONIAL DE LA LEY GENERAL DE CONTABILIDAD GUBERNAMENTAL",
				props.Text{Size: 8, Style: consts.Bold})
		})
	})
	m.Row(10, func() {
		m.Col(5, func() {
			m.Signature(doc.SignaturesInfo.Director.FullName, props.Font{Color: color.NewBlack()})
		})
		m.ColSpace(2)
		m.Col(5, func() {
			m.Signature("ING JOSE ANGEL GUZMAN DAVALOS", props.Font{Color: color.NewBlack()})
		})
	})
	m.Row(5, func() {
		m.Col(5, func() {
			m.Text("DIRECTOR DEL CENTRO ESTATAL DE CANCEROLOGIA", props.Text{Size: 8, Top: 10, Align: consts.Center})
		})
		m.ColSpace(2)
		m.Col(5, func() {
			m.Text("ENCARGADO DE ACTIVO FIJO", props.Text{Size: 8, Top: 10, Align: consts.Center})
		})
	})
	m.Line(10, props.Line{Color: color.NewWhite()})
	m.Row(6, func() {
		m.Col(5, func() {
			m.Signature(doc.SignaturesInfo.Administrator.FullName, props.Font{Color: color.NewBlack()})
		})
		m.ColSpace(2)
		m.Col(5, func() {
			m.Signature(doc.FixedAssetRequest.Department.ResponsibleUser.FullName, props.Font{Color: color.NewBlack()})
		})
	})
	m.Row(5, func() {
		m.Col(5, func() {
			m.Text("SUBDIRECTOR ADMINISTRATIVO DEL CENTRO ESTATAL DE CANCEROLOGIA", props.Text{Size: 8, Top: 10, Align: consts.Center})
		})
		m.ColSpace(2)
		m.Col(5, func() {
			m.Text("RESPONSABLE DEL AREA", props.Text{Size: 8, Top: 10, Align: consts.Center})
		})
	})
}
