package common

// Department represents an organizational unit in the ERP.
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Employee belongs to a department and can be used by salary/project services.
type Employee struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DepartmentID string `json:"departmentId"`
	Role         string `json:"role"`
}

// Salary stores compensation details for an employee.
type Salary struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employeeId"`
	Monthly    float64 `json:"monthly"`
	Currency   string  `json:"currency"`
}

// Project tracks delivery timeline for initiatives in ERP.
type Project struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Department  string   `json:"department"`
	Members     []string `json:"members"`
	Timeline    string   `json:"timeline"`
	StartDate   string   `json:"startDate"`
	ExpectedEnd string   `json:"expectedEnd"`
}
