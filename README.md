# argument

Get program arguments easily!


## Example

Add argument(s):

```go
import "github.com/oklookat/argument"

func main() {
    var argumentus = argument.New()

    argumentus.Add("username", "u", func(values []string) {
		
	})

    argumentus.Start()
}
```

Run your program with argument(s):

```
./myProgram --username

./myProgram --username=john

./myProgram --username john alice paul

./myProgram -u john alice paul
```


In values slice you got values like this:

```
nil

[john]

[john, alice, paul]
```
