Kombinasi Slice & Struct
Slice dan struct bisa dikombinasikan seperti pada slice dan map , caranya pun mirip,
cukup tambahkan tanda [] sebelum tipe data pada saat deklarasi.

type person struct {
name string
age int
}

var allStudents = []person{
{name: "Wick", age: 23},
{name: "Ethan", age: 23},
{name: "Bourne", age: 22},
}

for _, student := range allStudents {
fmt.Println(student.name, "age is", student.age)
}

-----------------------------------langsung
var allStudents = []struct {
	person
	grade int
}{
	{person: person{"wick", 21}, grade: 2},
	{person: person{"ethan", 22}, grade: 3},
	{person: person{"bond", 21}, grade: 3},
}

for _, student := range allStudents {
	fmt.Println(student)
}
