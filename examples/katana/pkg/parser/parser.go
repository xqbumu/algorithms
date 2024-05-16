package parser

import (
	"errors"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

type Row struct {
	Date   string
	Status string
	City   string
}

func (r *Row) recheck() error {
	date, err := time.Parse("02.01.2006", r.Date)
	if err != nil {
		return err
	}
	r.Date = date.Format("2006-01-02")
	return nil
}

func CdekExpress(p *rod.Page) ([]Row, error) {
	status := p.MustElement(".office-page-track__form-information-status")
	if !strings.Contains(status.MustText(), "Status") {
		return nil, errors.New("can not find status")
	}

	details := p.MustElements(".office-page-track__form-detail-items > .office-page-track__form-detail-item ")

	rows := make([]Row, 0, 10)
	for _, detail := range details {
		row := Row{
			Date:   detail.MustElement(".office-page-track__form-detail-date").MustText(),
			Status: detail.MustElement(".office-page-track__form-detail-status").MustText(),
			City:   detail.MustElement(".office-page-track__form-detail-city").MustText(),
		}
		if err := row.recheck(); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func Sibtrans(p *rod.Page) ([]Row, error) {
	result := p.MustElement(".status-result > tbody")
	details := result.MustElements("tr")
	if len(details) == 0 {
		return nil, errors.New("can not find result")
	}

	rows := make([]Row, 0, 10)
	for idx, detail := range details {
		var eles rod.Elements
		if idx == 0 {
			eles = detail.MustElements("th")
		} else {
			eles = detail.MustElements("td")
		}
		row := Row{
			Date:   eles[1].MustText(),
			Status: eles[2].MustText(),
			City:   eles[2].MustText(),
		}
		if err := row.recheck(); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}
