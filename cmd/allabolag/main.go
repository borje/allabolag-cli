package main

import (
	"flag"
	"log"

	"github.com/vertan/allabolag-cli/output"
	"github.com/vertan/allabolag-cli/scrape"
)

const minPositionalArgs = 1

func main() {
	// Parse flags
	terse := flag.Bool("t", false, "print company information in terse form")
	asJSON := flag.Bool("json", false, "print company information as JSON")
	flag.Parse()

	// Search term is a required argument
	if flag.NArg() < minPositionalArgs {
		flag.Usage()
		log.Fatal("missing required argument: search term")
	}

	query := flag.Arg(0)
	scraper := scrape.NewAllaBolagScraper()

	run(scraper, query, *terse, *asJSON)
}

func run(s scrape.CompanyInfoScraper, query string, terse bool, asJSON bool) {
	companies := getCompanies(s, query)
	if len(companies) == 0 {
		output.PrintNoResult(query)
		return
	}

	details := getDetails(s, companies[0])
	if details == nil {
		output.PrintNoResult(query)
		return
	}

	outputDetails(*details, terse, asJSON)
}

func getCompanies(s scrape.CompanyInfoScraper, query string) []scrape.Company {
	companies, _ := s.Search(query)
	// TODO: Handle parsing failure

	return companies
}

func getDetails(s scrape.CompanyInfoScraper, company scrape.Company) *scrape.CompanyDetails {
	details, _ := s.Details(company)
	// TODO: Handle parsing failure

	return details
}

func outputDetails(details scrape.CompanyDetails, terse bool, asJSON bool) {
	if asJSON {
		output.PrintJSON(details)
		return
	}
	if terse {
		output.PrintTerse(details)
		return
	}
	output.PrintSummary(details)
}
