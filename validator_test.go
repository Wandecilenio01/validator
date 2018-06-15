package validator

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

// WARNING: url case Test is not working for tests like www.google.com, but it works for http://www.example.com/index.html

// MyModel - Tests struct
type MyModel struct {
	ID             int64          `json:"id" struct-validator:"min:3|max:20"`
	Name           string         `json:"name" struct-validator:"regex:^[0-9]*$|required"`
	Age            int64          `json:"age" struct-validator:"min:3|max:20"`
	CreateAt       time.Time      `json:"createAt" struct-validator:"after_or_equal_date:today+3"`
	Email          string         `json:"email" struct-validator:"email|required_without_all:Site,JSON"`
	Site           string         `json:"site" struct-validator:"url"`
	IPv4           string         `json:"ipv4" struct-validator:"ipv4"`
	JSON           string         `json:"json" struct-validator:"json"`
	AlphaDashField string         `json:"alphaDash" struct-validator:"alpha_dash_space"`
	AlphaNumField  string         `json:"alphaNum" struct-validator:"alpha_num_space"`
	MyIntArray     []int          `json:"MyIntArray" struct-validator:"required_without_all:MyFloat32Array,MyUintptrArray"`
	MyFloat32Array []float32      `json:"MyFloat32Array" struct-validator:""`
	MyUintptrArray []uintptr      `json:"MyUintptrArray" struct-validator:""`
	MyAnotherModel MyAnotherModel `json:"MyAnotherModel" struct-validator:""`
	MyInt          int
	MyBool         bool `json:"MyBool" struct-validator:""`
}

type MyAnotherModel struct {
	FieldOne string
	FieldTwo int
}

var (
	examples = []MyModel{
		// Same as example.go
		{2, "", 21, time.Now().AddDate(0, 0, +3), "as", "", "as", "", "&&&**%%///\\\\%s", "&&&**%%///\\\\", []int{3}, nil, nil, *new(MyAnotherModel), 0, true},
		// All variables with no problem
		{3, "1900", 18, time.Now().AddDate(0, 0, +4), "robert@gmail.com", "http://www.example.com/index.html", "192.0.2.1", `{"name":"Robert"}`, "Joseph - it_administrator - 764", "Joseph764", []int{3}, []float32{1}, []uintptr{2}, *new(MyAnotherModel), 10, true},
		// All variables "empty"
		{0, "", 0, time.Time{}, "", "", "", "", "", "", []int{}, []float32{}, []uintptr{}, *new(MyAnotherModel), 0, false},
		// All variables "empty" except email. required_without_all test.
		{0, "", 0, time.Time{}, "example@hotmail.com", "", "", "", "", "", []int{}, []float32{}, []uintptr{}, *new(MyAnotherModel), 0, false},
		// ID and Age over max
		{40, "1900", 21, time.Now().AddDate(0, 0, +4), "robert@gmail.com", "http://www.example.com/index.html", "192.0.2.1", `{"name":"Robert"}`, "Joseph - it_administrator - 764", "Joseph764", []int{3}, []float32{1}, []uintptr{2}, *new(MyAnotherModel), 10, true},
	}
	stringsTest = [][]string{
		// Parameters not in a row
		[]string{"name", "id", "date"},
		// Only json
		[]string{"json"},
		// required_without_all function test
		[]string{"email"},
		// Only age
		[]string{"age"},
		// All parameters that have validations
		[]string{"id", "name", "age", "createAt", "email", "site", "ipv4", "json", "alphaDash", "alphaNUm", "MyIntArray"},
		// Only age and ID
		[]string{"age", "id"},
	}
	errorsTest = [][]error{
		[]error{errors.New("The ID cannot be less than 3, the value informed was 2."), errors.New(`The Name cannot have length less than 1, the informed value was "".`)},
		[]error{errors.New("The Age cannot be greater than 20, the value informed was 21.")},
		[]error{errors.New("The Email is not a valid required_without_all, because if all fields: (Site,JSON) are not filled, then Email needs to be filled.")},
		[]error{errors.New("The ID cannot be greater than 20, the value informed was 40."), errors.New("The Age is over max value.")},
		[]error{errors.New("The ID cannot be less than 3, the value informed was 2."), errors.New(`The Name cannot have length less than 1, the informed value was "".`), errors.New("The Age cannot be greater than 20, the value informed was 21."), errors.New(`The Email is not a valid email, the informed value was "as".`), errors.New(`The IPv4 is not a valid ipv4, the informed value was "as".`), errors.New(`The AlphaDashField is not a valid alpha_dash_space, the informed value was "&&&**%%///\\%s".`), errors.New(`The AlphaNumField is not a valid alpha_num_space, the informed value was "&&&**%%///\\".`)},
		[]error{errors.New("The ID cannot be less than 3, the value informed was 0."), errors.New(`The Name cannot have length less than 1, the informed value was "".`), errors.New("The Age cannot be less than 3, the value informed was 0."), errors.New("The CreateAt have to be after or equals to 2018-6-18, the date informed was 0001-1-1."), errors.New("The Email is not a valid required_without_all, because if all fields: (Site,JSON) are not filled, then Email needs to be filled."), errors.New("The MyIntArray is not a valid required_without_all, because if all fields: (MyFloat32Array,MyUintptrArray) are not filled, then MyIntArray needs to be filled.")},
		[]error{errors.New("The ID cannot be greater than 20, the value informed was 40."), errors.New("Invalid name.")},
	}
	messagesTest = map[string]map[string]string{
		"*": map[string]string{
			"min": "The {{.fieldName}} is under min value.",
		},
		"Age": map[string]string{
			"max": "The {{.fieldName}} is over max value.",
		},
	}
)

func TestValidate(t *testing.T) {
	t.Log("\nIt tests all fields in struct with struct-validator rules\n")
	if !reflect.DeepEqual(Validate(examples[0], nil), errorsTest[4]) {
		t.Log("\nTests the example.go\n")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
	if Validate(examples[1], nil) != nil {
		errorsReceived := Validate(examples[1], nil)
		t.Log("\nTests with correct values\n")
		t.Errorf("\nReceived: %v.\nShould be: nil.\n", errorsReceived)
	}
	if !reflect.DeepEqual(Validate(examples[2], nil), errorsTest[5]) {
		t.Log("Tests with default zero values")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
}

func TestValidateFields(t *testing.T) {
	t.Log("\nIt tests if validator tests only the string array sended\n")
	if !reflect.DeepEqual(ValidateFields(examples[0], stringsTest[0], nil), errorsTest[0]) {
		t.Log("\nTests two incorrect fields and one correct\n")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
	if ValidateFields(examples[1], stringsTest[4], nil) != nil {
		errorsReceived := ValidateFields(examples[1], stringsTest[4], nil)
		t.Log("\nTests all parameters with validations and no errors\n")
		t.Errorf("\nReceived: %v.\nShould be: nil.\n", errorsReceived)
	}
	if ValidateFields(examples[1], stringsTest[1], nil) != nil {
		errorsReceived := ValidateFields(examples[1], stringsTest[1], nil)
		t.Log("\nTests json with no error\n")
		t.Errorf("\nReceived: %v.\nShould be: nil.\n", errorsReceived)
	}
	if !reflect.DeepEqual(ValidateFields(examples[4], stringsTest[3], nil), errorsTest[1]) {
		t.Log("\nTests age over 20\n")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
	if ValidateFields(examples[3], stringsTest[2], nil) != nil {
		errorsReceived := ValidateFields(examples[3], stringsTest[2], nil)
		t.Log("\nTests required_without_all with no error\n")
		t.Errorf("\nReceived: %v.\nShould be: nil.\n", errorsReceived)
	}
	if !reflect.DeepEqual(ValidateFields(examples[2], stringsTest[2], nil), errorsTest[2]) {
		t.Log("\nTests required_without_all with error\n")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
	if !reflect.DeepEqual(ValidateFields(examples[4], stringsTest[5], messagesTest), errorsTest[3]) {
		t.Log("\nTests age over 20 with custom message\n")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
}

func TestAddCustomValidator(t *testing.T) {
	type CustomValidatorModel struct {
		ID   int64  `json:"id" struct-validator:"min:3|max:20"`
		Name string `json:"name" struct-validator:"name"`
	}
	testModel := CustomValidatorModel{40, "Err"}
	t.Log("\nIt tests if validator tests only the custom validator sended\n")
	if err := AddCustomValidator("string", "name", func(messageInput MessageInput) error {
		if len(messageInput.FieldValue.(string)) >= 4 {
			return nil
		}
		return errors.New("Invalid name.")
	}); err != nil {
		t.Errorf("\nCannot add custom validator.\n")
	}
	if !(reflect.DeepEqual(Validate(testModel, nil), errorsTest[6])) {
		t.Log("\nIt tests if custom validator in name is working\n")
		t.Errorf("\nReceived: false.\nShould be: true.\n")
	}
	testModel = CustomValidatorModel{18, "Name long enough."}
	if Validate(testModel, nil) != nil {
		errorsReceived := Validate(testModel, nil)
		t.Log("\nIt tests with correct data\n")
		t.Errorf("\nReceived: %v.\nShould be: nil.\n", errorsReceived)
	}
}
