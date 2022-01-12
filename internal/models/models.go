package models

import (
	"net/http"
	"sort"
	"strconv"
	"zimpler/internal/crawler"
)

type Models struct {
}

func NewModels() Models {
	return Models{}
}

type Customer struct {
	Name  string `json:"name"`
	Candy string `json:"candy"`
	Eaten int    `json:"eaten"`
}

type TopCustomers struct {
	Name           string `json:"name"`
	FavouriteSnack string `json:"favouriteSnack"`
	TotalSnacks    int    `json:"totalSnacks"`
}

func (m *Models) ExtractCustomers() ([]*Customer, error) {
	r, err := http.Get("https://candystore.zimpler.net/")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	root, err := crawler.NewCrawler(r.Body)
	if err != nil {
		return nil, err
	}

	table, err := root.FindElementByID("top.customers")
	if err != nil {
		return nil, err
	}

	tbody, err := table.FindElementByTag("tbody")
	if err != nil {
		return nil, err
	}

	children, err := tbody.Children()
	if err != nil {
		return nil, err
	}

	var customers []*Customer
	for _, tr := range children {
		tds, err := tr.Children()
		if err != nil {
			return nil, err
		}

		var c Customer
		for i, td := range tds {
			data, _ := td.InnerText()

			switch i {
			case 0:
				c.Name = data
			case 1:
				c.Candy = data
			case 2:
				c.Eaten, err = strconv.Atoi(data)
				if err != nil {
					return nil, err
				}
			default:
			}
		}
		customers = append(customers, &c)
	}

	return customers, nil
}

func (m *Models) GetTopCustomers() ([]*TopCustomers, error) {
	customers, err := m.ExtractCustomers()
	if err != nil {
		return nil, err
	}

	totals := make(map[string]int)
	favorites := make(map[string]map[string]int)

	for _, c := range customers {
		if _, ok := totals[c.Name]; !ok {
			totals[c.Name] = 0
		}

		totals[c.Name] += c.Eaten

		if _, ok := favorites[c.Name]; !ok {
			favorites[c.Name] = make(map[string]int)
		}

		if _, ok := favorites[c.Name][c.Candy]; !ok {
			favorites[c.Name][c.Candy] = 0
		}

		favorites[c.Name][c.Candy] += c.Eaten
	}

	var topCustomers []*TopCustomers
	for name, total := range totals {
		var tc TopCustomers
		tc.Name = name
		tc.TotalSnacks = total

		var fav string
		for snack, eaten := range favorites[name] {
			if fav == "" {
				fav = snack
				continue
			}

			if favorites[name][fav] < eaten {
				fav = snack
			}
		}
		tc.FavouriteSnack = fav

		topCustomers = append(topCustomers, &tc)
	}

	sort.SliceStable(topCustomers, func(i, j int) bool {
		return topCustomers[i].TotalSnacks > topCustomers[j].TotalSnacks
	})

	return topCustomers, nil
}
