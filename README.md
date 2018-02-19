Golang Validator 
=============
* [Usage](#usage)
* [Validator Key Types](#validator-key-types)
    * [Rules](#rules)
    * [Types](#types)

A GoLang validator to validate structs.

## Usage

Import this package typing:

    go get -u github.com/Wandecilenio01/validator

Define your model:
```
type MyModel struct {
    ID   int64  `json:"id" struct-validator:"min:3|max:20"`
    Name string `json:"name"`
    Age  int64  `json:"age"     struct-validator:"min:3|max:20"`
}
```
Then execute the validator with your struct:
```
func main() {
    onePerson := MyModel{1, "NameofPerson", 21}
    errors := validator.Validate(onePerson, nil)
    if len(errors) > 0 {
        for eindex, err := range errors {
            fmt.Println("Error Nº", eindex, "Error ->", err.Error())
        }
    }
}
```
You can add your custom validator if you want, you just need to type the following code:
```
validator.AddCustomValidator("string", "name", func(messageInput validator.MessageInput) error {
    if len(messageInput.FieldValue.(string)) < 4 {
        return fmt.Errorf("Erro: My err ...")
    }
    return nil
})
```
At that line, ```validator.AddCustomValidator("string", "name" ...```, ```"string"``` is the **[Validator Key Type](#validator-key-types)** and "name" is the **[Rule](#rules)**, more details about that terms, see the links.

## Validator Key Types

The validators are the data types used by the rules.

### Rules

#### Types

* **numeric**: Represents types int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32 and float64. Rules:
    * **min**:
    * **max**:
* **string**: Represents the string type.

