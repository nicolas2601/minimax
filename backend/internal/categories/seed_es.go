package categories

// DefaultCategoriesES returns the seed list of categories for a new user.
// Onboarding creates these on first authenticated request to /categories/seed.
var DefaultCategoriesES = []Category{
	// Gastos (expense)
	{Name: "Alimentación", Type: TypeExpense, Icon: ptr("utensils"), Color: ptr("#f4c5a8")},
	{Name: "Transporte", Type: TypeExpense, Icon: ptr("car"), Color: ptr("#a8c8e8")},
	{Name: "Vivienda", Type: TypeExpense, Icon: ptr("home"), Color: ptr("#a7e5d3")},
	{Name: "Salud", Type: TypeExpense, Icon: ptr("heart"), Color: ptr("#e8b8c4")},
	{Name: "Entretenimiento", Type: TypeExpense, Icon: ptr("film"), Color: ptr("#c8b8e0")},
	{Name: "Compras", Type: TypeExpense, Icon: ptr("shopping-bag"), Color: ptr("#a8c8e8")},
	{Name: "Servicios", Type: TypeExpense, Icon: ptr("zap"), Color: ptr("#f4c5a8")},
	{Name: "Educación", Type: TypeExpense, Icon: ptr("book"), Color: ptr("#a7e5d3")},
	{Name: "Otros gastos", Type: TypeExpense, Icon: ptr("more-horizontal"), Color: ptr("#a8a29e")},

	// Ingresos (income)
	{Name: "Salario", Type: TypeIncome, Icon: ptr("briefcase"), Color: ptr("#16a34a")},
	{Name: "Freelance", Type: TypeIncome, Icon: ptr("laptop"), Color: ptr("#16a34a")},
	{Name: "Inversiones", Type: TypeIncome, Icon: ptr("trending-up"), Color: ptr("#16a34a")},
	{Name: "Otros ingresos", Type: TypeIncome, Icon: ptr("plus-circle"), Color: ptr("#16a34a")},
}

func ptr[T any](v T) *T { return &v }