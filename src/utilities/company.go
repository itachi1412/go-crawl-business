package utilities

import (
	"net/http"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/mapstructure"
)

type Company struct {
	Name                 string `json:"mame"`
	AlternateName        string `json:"alternate_name"`
	TaxId                string `json:"tax_id"`
	Status               string `json:"status"`
	TaxRegisterAddress   string `json:"tax_register_address"`
	Address              string `json:"address"`
	Phone                string `json:"phone"`
	Founder              string `json:"founder"`
	President            string `json:"president"`
	LicenseDate          string `json:"license_date"`
	StartDate            string `json:"start_date"`
	ReceiveAccountDate   string `json:"receive_account_date"`
	FinancialYear        string `json:"financial_year"`
	Employees            string `json:"employees"`
	LevelChapterItemType string `json:"level_chapter_item_type"`
	BankAccount          string `json:"bank_account"`
}

type Companies struct {
	TotalCompanies int       `json:"total_companies"`
	List           []Company `json:"companies"`
}

func NewCompanies() *Companies {
	return &Companies{}
}

func formatField(a string) string {
	a = strings.Replace(a, "Tên doanh nghiệp:", "Name", 1)
	a = strings.Replace(a, "Tên giao dịch", "AlternateName", 1)
	a = strings.Replace(a, "Mã số thuế:", "TaxId", 1)
	a = strings.Replace(a, "Tình trạng hoạt động:", "Status", 1)
	a = strings.Replace(a, "Nơi đăng ký quản lý:", "TaxRegisterAddress", 1)
	a = strings.Replace(a, "Địa chỉ:", "Address", 1)
	a = strings.Replace(a, "Điện thoại:", "Phone", 1)
	a = strings.Replace(a, "Đại diện pháp luật:", "Founder", 1)
	a = strings.Replace(a, "Giám đốc:", "President", 1)
	a = strings.Replace(a, "Ngày cấp giấy phép:", "LicenseDate", 1)
	a = strings.Replace(a, "Ngày bắt đầu hoạt động:", "StartDate", 1)
	a = strings.Replace(a, "Ngày nhận TK:", "ReceiveAccountDate", 1)
	a = strings.Replace(a, "Năm tài chính:", "FinancialYear", 1)
	a = strings.Replace(a, "Số lao động:", "Employees", 1)
	a = strings.Replace(a, "Cấp Chương Loại Khoản:", "LevelChapterItemType", 1)
	a = strings.Replace(a, "TK ngân hàng:", "BankAccount", 1)
	return a
}

func (companies *Companies) ExtractInfomation(url string, client *http.Client) error {
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"user-agent": {browser.Firefox()},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	var s string
	doc.Find("div.company-info-section").First().Each(func(i int, sectionHtml *goquery.Selection) {
		sectionHtml.Find("div.responsive-table-cell").Each(func(i int, cellHtml *goquery.Selection) {
			s += strings.TrimSpace(cellHtml.Text()) + "\n"
		})
	})

	s = formatField(s)
	x := strings.Split(s, "\n")
	y := make(map[string]string)
	for i := 0; i < len(x)-1; i += 2 {
		y[x[i]] = x[i+1]
	}
	Company := Company{}
	mapstructure.Decode(y, &Company)

	companies.TotalCompanies++
	companies.List = append(companies.List, Company)

	return nil
}
