Golang Validator
=============

* [Usage](#usage)
* [Validator Key Types](#validator-key-types)
    * [Rules](#rules)
    * [Types](#types)
* [Custom Validations](#custom-validations)
* [Custom Messages](#custom-messages)
* [Message Input](#message-input)
* [Validate Custom Fields](#validate-custom-fields)
* [Set Tag Name](#set-tag-name)

A GoLang validator to validate structs.

## Usage

Import this package typing:

    go get -u github.com/Wandecilenio01/validator

Define your model:

```Golang
type MyModel struct {
    ID   int64  `json:"id" struct-validator:"min:3|max:20"`
    Name string `json:"name"`
    Age  int64  `json:"age" struct-validator:"min:3|max:20"`
}
```

Then execute the validator with your struct:

```Golang
func main() {
    onePerson := MyModel{1, "NameofPerson", 21}
    errors := validator.Validate(onePerson, nil)
    if len(errors) > 0 {
        for _, err := range errors {
            fmt.Println("Error -> ", err.Error())
        }
    }
}
```

You can add your custom validator if you want, you just need to type the following code:

```Golang
validator.AddCustomValidator("string", "name", func(messageInput validator.MessageInput) error {
    if len(messageInput.FieldValue.(string)) < 4 {
        return fmt.Errorf("Erro: My err ...")
    }
    return nil
})
```

At that line, ```validator.AddCustomValidator("string", "name" ...```, ```"string"``` is the **[Validator Key Type](#validator-key-types)** and ```"name"``` is the **[Rule](#rules)**, more details about that terms, click on links. More details about rules can be found in **[Types](#types)**. At the line ```messageInput.FieldValue.(string)``` we have a casting, to get the field value from a **[Message Input](#message-input)**.

## Validator Key Types

A 'validator key type' is a name used by the golang-validator to reference rules and golang type data.

### Rules

Rules are names given to handler of each type.

#### Types

Below there's a list of validator key type with sublists of rules.

* **numeric**: Represents types int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32 and float64. Rules:
    * **min**: Minimum value acceptable by field, ```(min:3)```;
    * **max**: Maximum value acceptable by field, ```(min:45)```.
* **string**: Represents the string type.
    * **min**: Minimum length acceptable by field, ```(min:3)```;
    * **max**: Maximum length acceptable by field, ```(min:65)```;
    * **length**: Exact length acceptable by field, ```(min:43)```;
    * **email**: The field value have to be a valid email;
    * **url**: The field value have to be a valid URL;
    * **ipv4**: The field value have to be a valid IPv4;
    <!-- * **IPv6**: The field value have to be a valid IPv6; -->
    * **json**: The field value have to be a valid JSON text;
    * **alpha**: The field value can only have Latin alphabet char's;
    * **alpha_num**: Rule *alpha* + numbers;
    * **alpha_dash**: Rule *alpha_num* + '-' + '_';
    * **alpha_space**: The field value can only have Latin alphabet char's and spaces;
    * **alpha_num_space**: Rule *alpha_space* + numbers;
    * **alpha_dash_space**: Rule *alpha_num_space* + '-' + '_';
    * **regex**: The field value have to match the specified regex, ```(regex:^[0-9]*$)``` or ```(regex:^((\\d{3}).(\\d{3}).(\\d{3})-(\\d{2}))*$)```;
    * **required**: The field value cannot be empty.
    * **required_with**: The field under validation must be present and not empty only if any of the other specified fields are not empty, ```(required_with:field1,field2)```;
	* **required_with_all**: The field under validation must be present and not empty only if all of the other specified fields are not empty, ```(required_with_all:field1,field2)```;
	* **required_without**: The field under validation must be present and not empty only when any of the other specified fields are empty, ```(required_without:field1,field2)```;
    * **required_without_all**: The field under validation must be present and not empty only when all of the other specified fields are empty, ```(required_without_all:field1,field2)```.
* **timestamp**: Represents the ```time.Time``` type.
    In this type the rule value used is ```today```, ```today+1``` represents tomorrow, ```today-1``` represents yesterday and so on, for example: ```today-2```, ```today+3```, ... .
    * **after**: The field value have to be after the specified time, ```(after:today)```;
    * **before**: The field value have to be before the specified time, ```(before:today)```;
    * **equal**: The field value have to be equals to the specified time, ```(equal:today)```;
    * **after_or_equal**: *after* or *equal*, ```(after_or_equal:today)```;
    * **before_or_equal**: *before* or *equal*, ```(before_or_equal:today)```;
    * **after_date**: The field value have to be after the specified date, considers only the date part of ```time.Time```, ```(after_date:today)```;
    * **before_date**: The field value have to be before the specified date, considers only the date part of ```time.Time```, ```(before_date:today)```;
    * **equal_date**: The field value have to be equal the specified date, considers only the date part of ```time.Time```, ```(equal_date:today)```;
    * **after_or_equal_date**: *after_date* or *equal_date*, ```(after_or_equal:today)```;
    * **before_or_equal_date**: *before_date* or *equal_date*, ```(before_or_equal:today)```.
* **arrray**: Represents the any array used, only arrays, not pointers.
    * **min**: Minimum length acceptable by array, ```(min:2)```.
    * **max**: Maximum length acceptable by array, ```(max:3)```.
    * **length**:  Exact length acceptable by array, ```(length:3)```.
    * **distinct**: The field array cannot have repeated values.
    * **required**: The field value cannot be an empty array or nil.
    * **required_with**: The field under validation must be present and not empty only if any of the other specified fields are not empty, ```(required_with:field1,field2)```;
    * **required_with_all**: The field under validation must be present and not empty only if all of the other specified fields are not empty, ```(required_with_all:field1,field2)```;
    * **required_without**: The field under validation must be present and not empty only when any of the other specified fields are empty, ```(required_without:field1,field2)```;
    * **required_without_all**: The field under validation must be present and not empty only when all of the other specified fields are empty, ```(required_without_all:field1,field2)```.

## Custom Validations

To add custom validations we need to define a code like bellow.

```Golang
validator.AddCustomValidator("string", "name", func(messageInput validator.MessageInput) error {
    if len(messageInput.FieldValue.(string)) >= 4 {
        return nil
    }
    return errors.New("Invalid name")
})
```

The line ```validator.AddCustomValidator("string", "name", ...``` we define the validator key type and the rule name for our handler (OBS: The golang-validator doesn't let you change the native validators presented in **[section](#validator-key-types)**).

## Message Input

The message input is the data structure used by golang-validator as input.

```Golang
type MessageInput struct {
    FieldName        string
    FieldType        reflect.Type
    FieldValue       interface{}
    ValidatorKeyType string
    RuleName         string
    RuleValue        string
    CustomMessages   map[string]map[string]string
}
```

* **FieldName**: Represents the attribute name of mapped struct. For example, in ```Name string `struct-validator:"required"` ```, the *FieldName* will be *Name*.
* **FieldType**: Represents the attribute type of mapped struct. For example, in ```Name string `struct-validator:"required"` ```, the *FieldType* will be a ```reflect.Type``` that represents a string type.
* **FieldValue**: Represents the attribute value of mapped struct. The value will be an interface, so the developer will responsible to do a cast to use the original value from this attribute.
* **ValidatorKeyType**: Represents the **[Validator Key Type](#validator-key-types)**.
* **RuleName**: Represents the rule used, more **[info](#validator-key-types)**.
* **RuleValue**: Represents the rule value used, for example, in ```Name string `struct-validator:"required"` ```, the rule value will be ```required```, more **[info](#validator-key-types)**.
* **CustomMessages**: Represents the map of messages, the first key represents the field name of mapped struct and the second key represents the rule name, more **[info](#custom-messages)**.

## Custom Messages

To define custom messages we need to define a variable like bellow.

```Golang
messages := map[string]map[string]string{
    "*": map[string]string{
        "min": "The min value for {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
    },
    "Age": map[string]string{
        "max": "The max {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
    },    
}
errors := validator.Validate(onePerson, messages)
```

At the line ```"*": map[string]string{ ...``` we are defining messages to every attribute, if you want to define a message to only attribute, you will do like code at ```"Age": map[string]string{...```. In the line ```"min": "The min value for {{.fieldName}} should be ...",``` is defined a message for the ```"min"``` rule, and the ```{{.fieldName}}``` represents the template text.

## Validate Custom Fields

To define custom fields to be validated.

Define your model:

```Golang
type MyModel struct {
    ID   int64  `json:"id" struct-validator:"min:3|max:20"`
    Name string `json:"name"`
    Age  int64  `json:"age" struct-validator:"min:3|max:20"`
}
```

Then execute the validator with your struct:

```Golang
func main() {
    onePerson := MyModel{1, "NameofPerson", 21}
	errors := validator.ValidateFields(onePerson, []string{"id"}, nil)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println("Error -> ", err.Error())
		}
	}
}
```

Only the fields passed by the string array will be validated (in this example, age will not trigger an error), follow the output of this example:

    Error ->  The ID cannot be less than 3, the value informed was 1.

## Set Tag Name

To define a custom tag name to substitute "struct-validator".

Define your own tag:

    validator.SetTag("validate")

Then define your model with your new tag:

```Golang
type MyModel struct {
    ID   int64  `json:"id" validate:"min:3|max:20"`
    Name string `json:"name"`
    Age  int64  `json:"age" validate:"min:3|max:20"`
}
```