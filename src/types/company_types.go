package types

type CompanyTypes string

type companyEnum struct {
	Client   CompanyTypes
	Catering CompanyTypes
}

// CompanyTypesEnum enum
var CompanyTypesEnum = companyEnum{
	Client:   "client",
	Catering: "catering",
}
