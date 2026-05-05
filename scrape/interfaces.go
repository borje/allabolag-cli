package scrape

// CompanyInfoScraper represents a generic scraper of company information.
type CompanyInfoScraper interface {
	Search(term string) ([]Company, error)
	SearchPersons(term string) ([]PersonResult, error)
	Details(c Company) (*CompanyDetails, error)
}

// PersonResult represents a person and their associated business activities.
type PersonResult struct {
	Name       string
	Age        int
	Location   string
	Businesses []Business
}

// Business represents a company a person has a role in.
type Business struct {
	Name  string
	Orgnr string
	Role  string
}

// Company represents a company.
type Company struct {
	Name     string
	Link     string
	Location string
}

// CompanyDetails represents further details about a company.
type CompanyDetails struct {
	Company
	Roles         []Role
	FiscalDetails []FiscalDetails
}

// Role represents a person or company with a role in a company.
type Role struct {
	Name  string
	Title string
	Group string
}

// FiscalDetails represents the financial information about a company for one year.
type FiscalDetails struct {
	Year    int
	Revenue int
	Result  int
}
