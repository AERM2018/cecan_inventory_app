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

type FixedAssetsReportDoc struct {
	FixedAssets []FixedAssetDetailed
	FromDate    time.Time
	ToDate      time.Time
	CreatedAt   time.Time
}

func (doc *FixedAssetsReportDoc) CreateDoc() (string, error) {
	cwd, _ := os.Getwd()
	doc.CreatedAt = time.Now()
	mto := pdf.NewMaroto(consts.Landscape, consts.Letter)
	// mto.SetPageMargins(12, 12, 12)
	common.InsertPdfHeader(mto, true)
	doc._addReportInfo(mto)
	doc._addFixedAssetsList(mto)
	outputFilePath := path.Join(cwd, "domain", "pdfs", "fixed_asset_report.pdf")
	err := mto.OutputFileAndClose(outputFilePath)
	if err != nil {
		return "", err
	}
	return outputFilePath, nil
}
func (doc FixedAssetsReportDoc) _addReportInfo(m pdf.Maroto) {
	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("REPORTE DE ACTIVO FIJO EN EL HOSPITAL", props.Text{Size: 15, Style: consts.Bold, Align: consts.Center})
		})
	})
	// m.Line(3, props.Line{Color: color.NewWhite()})
}

func (doc FixedAssetsReportDoc) _addFixedAssetsList(m pdf.Maroto) {

	headers := []string{"ETIQUETA (CLAVE DE INVENTA- RIO)", "DESCRIPCION", "MARCA Y MODELO", "SERIE", "TIPO", "ESTADO FISICO", "DEPTO.", "RESP. DE DEPTO", "DIR.", "ADMIN.", "FECHA DE INGRESO AL AREA"}
	gridSizes := []uint{1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	contents := make([][]string, 0)
	for _, fixedAsset := range doc.FixedAssets {
		fixedAssetInfo := []string{
			fixedAsset.Key,
			fixedAsset.Description,
			fmt.Sprintf("%v (%v)", fixedAsset.Brand, fixedAsset.Model),
			fixedAsset.Series,
			fixedAsset.Type,
			fixedAsset.PhysicState,
			fmt.Sprintf("%v (%v)", fixedAsset.DepartmentName, fixedAsset.DepartmentFloorNumber),
			fixedAsset.DepartmentResponsibleUserName,
			fixedAsset.DirectorUserName,
			fixedAsset.AdministratorUserName,
			fixedAsset.CreatedAt.Format("02/01/2006"),
		}
		contents = append(contents, fixedAssetInfo)
	}
	if len(contents) > 0 {
		m.Line(10, props.Line{Color: color.NewWhite()})
		m.SetBorder(true)
		m.Row(16, func() {
			for gzPos, h := range headers {
				m.Col(gridSizes[gzPos], func() {
					m.Text(h, props.Text{Left: 1.5, Top: 2, Extrapolate: false, Style: consts.Bold, Size: 9})
				})

			}
		})
		for _, row := range contents {
			m.Row(13, func() {
				m.SetBorder(true)
				for gzPos, col := range row {
					m.Col(gridSizes[gzPos], func() {
						m.Text(col, props.Text{Left: 1.5, Right: 1.5, Top: 1, Extrapolate: false, Size: 8})
					})
				}
				m.SetBorder(false)
			})
		}

	}
}
