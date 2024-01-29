Functionalities:
- reads mathematical expressions from an input file, evaluates and prints them

![screen cast](live_evaluator.gif)

Uses:
- fsnotify to watch the file changes
- each read reads from last offset for the new content
- github.com/maja42/goval is used for evaluation of the expressions
- evaluation is done on workers (go routines)


Usage
```
go run main.go

#open test.txt and start typing mathematical expressions and enter save
```