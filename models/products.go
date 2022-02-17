package models

import "GolangStore/database"

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float64
	Quantity    int
}

func ProductFinder() []Product {
	db := database.Connect()
	defer db.Close()
	selectAll, err := db.Query("select * from products")
	if err != nil {
		panic(err.Error())
	}
	p := Product{}
	products := []Product{}
	for selectAll.Next() {
		err = selectAll.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Quantity)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, p)
	}
	return products
}
