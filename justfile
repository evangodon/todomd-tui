dev:  
  go run . --file tmp/todo-example.md interactive

test:
  gotest -v ./...

todo:
  todomd interactive 
