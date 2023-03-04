# argp module

## introduction

`argp` is a dead simple module to build a cli tool using golang. Differs from heavy alternatives such as cobra, it is quite lightweight and enough to use. It use go generic to make
passing command options more easy.

Because `argp` is dead simple, it has limitations, for examples, it doesn't support list arguments, and parameters between parent command and sub command should not be the same.

`argp` use built-in `flag` modules to parse arguments, once arguments have been parsed, and the command has sub command, then the rest args will be passed to sub-command.

## How to use?

Using  `argp` is simple, and have only one way. First defines struct to save your options:

```go
type Option struct{
    Name string `arg:"name,the name of person"`
    Age int `arg:"age,the age of person"`
}
```
make sure to write down tag, this is the must do to bind it to flag set. The first part of tag is the argument name, and the second is usage.

Next step, defines Command object:

```go
var cmd := &Cmd[Option]{
    CmdName: "person",
    Usage: "read a person's name and age",
    Options: &Option{
        Age: 18,
    },
    Command: func(opt *Option) error {
        fmt.Printf("your name is %s, and you are %d years old!\n", opt.Name, opt.Age)
    }
}
```

To run the command, just use `cmd.Run(os.Args[1:])` to pass in the root arguments, note that for `Cmd` implentaions, don't pass in the name of the command itself.

If you want to use sub command, and want to use arguments from parent command, you should defines the option object first, and pass the pointer to parent's `Options` filed. 
You can visit your parent's option in the sub-command's `Command` function. Note that all flags for parent arguments should be geven before sub command start.
