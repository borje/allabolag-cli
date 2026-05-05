package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vertan/allabolag-cli/scrape"
)

// PrintTerse outputs company details in a terse format to the terminal.
func PrintTerse(c scrape.CompanyDetails) {
	fmt.Printf("%s\n", c.Company.Name)
	fmt.Printf("%s\n", c.Company.Link)

	if len(c.FiscalDetails) > 0 {
		fmt.Printf("Revenue (%d): %dk\n", c.FiscalDetails[0].Year, c.FiscalDetails[0].Revenue)
		fmt.Printf("Results (%d): %dk\n", c.FiscalDetails[0].Year, c.FiscalDetails[0].Result)
	}
}

// PrintSummary outputs company details in a summary format to the terminal.
func PrintSummary(c scrape.CompanyDetails) {
	fmt.Printf("%s\n", c.Company.Name)
	if c.Company.Location != "" {
		fmt.Printf("%s\n", c.Company.Location)
	}
	fmt.Printf("%s\n", c.Company.Link)
	if len(c.Roles) > 0 {
		fmt.Println("--------------------")
		printRolesTable(c.Roles)
	}
	if len(c.FiscalDetails) > 0 {
		fmt.Println("--------------------")
		printFiscalTable(c.FiscalDetails)
	}
}

// PrintJSON outputs company details as JSON.
func PrintJSON(c scrape.CompanyDetails) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(c)
}

// PrintPersonResultsJSON outputs person results as JSON.
func PrintPersonResultsJSON(persons []scrape.PersonResult) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(persons)
}

// PrintPersonResults outputs a list of persons and their business associations.
func PrintPersonResults(persons []scrape.PersonResult) {
	for i, p := range persons {
		if i > 0 {
			fmt.Println()
		}
		if p.Location != "" {
			fmt.Printf("%s (%d år) - %s\n", p.Name, p.Age, p.Location)
		} else {
			fmt.Printf("%s (%d år)\n", p.Name, p.Age)
		}
		for _, b := range p.Businesses {
			fmt.Printf("  %-40s %s  %s\n", b.Name, b.Orgnr, b.Role)
		}
	}
}

// PrintNoResult outputs a string for when there's no results..
func PrintNoResult(t string) {
	fmt.Printf("No result found for search term %s\n", t)
}

// printRolesTable renders tabular role data to the terminal.
func printRolesTable(roles []scrape.Role) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tTitle\tGroup")
	for _, r := range roles {
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, r.Title, r.Group)
	}
	w.Flush()
}

// printFiscalTable renders tabular financial data to the terminal.
func printFiscalTable(details []scrape.FiscalDetails) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Year\tRevenue\tResult")
	for _, v := range details {
		fmt.Fprintln(w, fmt.Sprintf("%d\t%dk\t%dk", v.Year, v.Revenue, v.Result))
	}
	w.Flush()
}
