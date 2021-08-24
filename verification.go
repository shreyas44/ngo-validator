package main

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func ScrapeNGOInfo(panId string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = chromedp.NewContext(ctx)

	var (
		panInputSel      = "#ctl00_SPWebPartManager1_g_d6877ff2_42a8_4804_8802_6d49230dae8a_txtPanAgency"
		searchButtonSel  = "#ctl00_SPWebPartManager1_g_d6877ff2_42a8_4804_8802_6d49230dae8a_btnSubmit"
		searchResultsSel = ".act_search_results"
		html             string
	)

	chromedp.Run(
		ctx,
		chromedp.Navigate(INCOME_TAX_URL),
		chromedp.WaitVisible(panInputSel),
		chromedp.SendKeys(panInputSel, panId),
		chromedp.Click(searchButtonSel),
		chromedp.WaitReady(searchResultsSel),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	if html == "" {
		return "", errors.New("unexpected error")
	}

	return html, nil
}

func GetNGOInfo(html string) (NGOInfo, bool) {
	var (
		info        = NGOInfo{}
		exists      = false
		dataIndices = map[int]string{
			0: "address",
			1: "state",
		}
	)

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find(".pan-id").Remove()
	info.Name = doc.Find(".faqsno-heading div").Text()

	if info.Name != "" {
		exists = true
	}

	doc.Find(".exempted-detail:first-child li").Each(func(i int, s *goquery.Selection) {
		if field, ok := dataIndices[i]; ok {
			switch field {
			case "address":
				info.Address = s.Find("span").Text()
			case "state":
				info.State = s.Find("span").Text()
			}
		}
	})

	return info, exists
}

func cleanName(name string) string {
	regex, _ := regexp.Compile("[.@,_]")
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = string(regex.ReplaceAll([]byte(name), []byte("")))
	return name
}

func VerifyNGO(input VerifyInput) error {
	html, err := ScrapeNGOInfo(input.PAN)
	if err != nil {
		return errors.New("unexpected_error")
	}

	info, exists := GetNGOInfo(html)

	switch {
	case !exists:
		return errors.New("invalid_pan")
	case cleanName(input.Name) != cleanName(info.Name):
		return errors.New("invalid_name")
	case !VerifyAddresses(input.Address, info.Address):
		return errors.New("invalid_address")
	}

	return nil
}

func VerifyAddresses(address1 string, address2 string) bool {
	return true
}
