package types

type companyEnum struct {
	Client   string
	Catering string
}

// CompanyTypesEnum enum
var CompanyTypesEnum = companyEnum{
	Client:   "client",
	Catering: "catering",
}
