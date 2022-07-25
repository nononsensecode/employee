package employee

type Employee struct {
	id   int64
	name string
	age  uint8
}

func New(name string, age uint8) Employee {
	return Employee{
		name: name,
		age:  age,
	}
}

func UnmarshalFromPersistence(id int64, name string, age uint8) Employee {
	return Employee{
		id:   id,
		name: name,
		age:  age,
	}
}

func (e Employee) ID() int64 {
	return e.id
}

func (e Employee) Name() string {
	return e.name
}

func (e Employee) Age() uint8 {
	return e.age
}
