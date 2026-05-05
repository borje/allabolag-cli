package scrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// AllaBolagScraper is an implementation of CompanyInfoScraper fetching info from allabolag.se.
type AllaBolagScraper struct{}

const (
	searchURL       = "https://www.allabolag.se/bransch-sök?q=%s"
	maxYearsToFetch = 5
)

type nextDataSearchCompanyList struct {
	Companies []struct {
		Name      string `json:"name"`
		CompanyID string `json:"companyId"`
		Orgnr     string `json:"orgnr"`
		Location  struct {
			Municipality string `json:"municipality"`
		} `json:"location"`
		Industries []struct {
			Name      string `json:"name"`
			CompanyID string `json:"companyId"`
		} `json:"industries"`
	} `json:"companies"`
}

type nextDataSearch struct {
	Props struct {
		PageProps struct {
			HydrationData struct {
				SearchStore struct {
					CompaniesByName nextDataSearchCompanyList `json:"companiesByName"`
					Companies       nextDataSearchCompanyList `json:"companies"`
				} `json:"searchStore"`
			} `json:"hydrationData"`
		} `json:"pageProps"`
	} `json:"props"`
}

type nextDataCompanyPage struct {
	Props struct {
		PageProps struct {
			Company struct {
				CompanyAccounts []struct {
					Year     string `json:"year"`
					Accounts []struct {
						Code   string `json:"code"`
						Amount string `json:"amount"`
					} `json:"accounts"`
				} `json:"companyAccounts"`
				Roles struct {
					RoleGroups []struct {
						Name  string `json:"name"`
						Roles []struct {
							Name string `json:"name"`
							Role string `json:"role"`
						} `json:"roles"`
					} `json:"roleGroups"`
				} `json:"roles"`
			} `json:"company"`
		} `json:"pageProps"`
	} `json:"props"`
}

func browserHeaders(r *colly.Request) {
	r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	r.Headers.Set("Accept-Language", "sv-SE,sv;q=0.9")
	r.Headers.Set("Sec-Fetch-Mode", "navigate")
	r.Headers.Set("Cookie", "i18next=sv")
}

// Search takes a search term as a parameter and searches allabolag.se for companies.
func (s *AllaBolagScraper) Search(term string) ([]Company, error) {
	c := colly.NewCollector()
	c.OnRequest(browserHeaders)

	var nd nextDataSearch
	idToLink := map[string]string{}

	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		json.Unmarshal([]byte(e.Text), &nd)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if !strings.HasPrefix(href, "/foretag/") {
			return
		}
		parts := strings.Split(strings.TrimRight(href, "/"), "/")
		companyID := parts[len(parts)-1]
		if _, exists := idToLink[companyID]; !exists {
			idToLink[companyID] = "https://www.allabolag.se" + href
		}
	})

	_ = c.Visit(fmt.Sprintf(searchURL, url.QueryEscape(term)))

	store := nd.Props.PageProps.HydrationData.SearchStore
	candidates := store.CompaniesByName.Companies
	if len(candidates) == 0 {
		candidates = store.Companies.Companies
	}

	companies := []Company{}
	for _, nc := range candidates {
		link, ok := idToLink[nc.CompanyID]
		if !ok {
			continue
		}
		companies = append(companies, Company{Name: nc.Name, Link: link, Location: nc.Location.Municipality})
	}
	return companies, nil
}

// Details returns details about a specific company.
func (s *AllaBolagScraper) Details(comp Company) (*CompanyDetails, error) {
	c := colly.NewCollector()
	c.OnRequest(browserHeaders)

	var nd nextDataCompanyPage

	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		json.Unmarshal([]byte(e.Text), &nd)
	})

	_ = c.Visit(comp.Link)

	accounts := nd.Props.PageProps.Company.CompanyAccounts
	if len(accounts) == 0 {
		return nil, errors.New("no fiscal data found")
	}

	limit := len(accounts)
	if limit > maxYearsToFetch {
		limit = maxYearsToFetch
	}

	fiscalDetails := []FiscalDetails{}
	for _, acc := range accounts[:limit] {
		year, err := strconv.Atoi(acc.Year)
		if err != nil {
			continue
		}
		var revenue, result int
		for _, a := range acc.Accounts {
			switch a.Code {
			case "SDI":
				revenue, _ = strconv.Atoi(a.Amount)
			case "ORS":
				result, _ = strconv.Atoi(a.Amount)
			}
		}
		fiscalDetails = append(fiscalDetails, FiscalDetails{Year: year, Revenue: revenue, Result: result})
	}

	roles := []Role{}
	for _, group := range nd.Props.PageProps.Company.Roles.RoleGroups {
		for _, r := range group.Roles {
			roles = append(roles, Role{Name: r.Name, Title: r.Role, Group: group.Name})
		}
	}

	return &CompanyDetails{Company: comp, Roles: roles, FiscalDetails: fiscalDetails}, nil
}

// NewAllaBolagScraper returns a new AllaBolagScraper
func NewAllaBolagScraper() *AllaBolagScraper {
	return &AllaBolagScraper{}
}
