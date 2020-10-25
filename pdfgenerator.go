package main

import (
	"bytes"
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"os"
	"time"
)

func generateQrsPdf(items []trackedItem) bytes.Buffer {
	begin := time.Now()

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	m.Row(15, func() {
		for i := 0; i < len(items); i = i + 3 {
			m.Col(6, func() {
				for j := 0; j < 3; j++ {
					if index := i + j; index < len(items) {
						item := items[index]
						m.QrCode(base+"/"+item.Id, props.Rect{
							Left:    5,
							Top:     5,
							Center:  false,
							Percent: 100,
						}) //items[0].Id)
					}
				}
			})
		}
	})

	//	err := m.OutputFileAndClose("example.pdf")
	pdfdata, err := m.Output()
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
	return pdfdata
}
