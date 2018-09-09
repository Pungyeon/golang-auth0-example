import { Component, OnInit } from '@angular/core';
import { TodoService, Todo } from '../service/todo.service';

@Component({
  selector: 'app-todo',
  templateUrl: './todo.component.html',
  styleUrls: ['./todo.component.css']
})
export class TodoComponent implements OnInit {

  activeTodos: Todo[];
  completedTodos: Todo[]
  todoTitle: string;
  todoDescription: string;

  constructor(private todoService: TodoService) { }

  ngOnInit() {
    this.getAll();
  }

  getAll() {
    this.todoService.getTodoList().subscribe((data: Todo[]) => {
      console.log(data);
      this.activeTodos = data.filter((a) => !a.complete);
      this.completedTodos = data.filter((a) => a.complete);
    });
  }

  addTodo() {
    var newTodo : Todo = {
      title: this.todoTitle,
      description: this.todoDescription,
      id: '',
      username: '',
      complete: false
    }
    this.todoService.addTodo(newTodo).subscribe(() => {
      this.getAll();
      this.todoTitle = '';
      this.todoDescription = '';
    });
  }

  completeTodo(todo: Todo) {
    this.todoService.completeTodo(todo).subscribe(() => {
      this.getAll();
    });
  }

  deleteTodo(todo: Todo) {
    this.todoService.deleteTodo(todo).subscribe(() => {
      this.getAll();
    })
  }
}
