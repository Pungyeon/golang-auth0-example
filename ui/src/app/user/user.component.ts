import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Router } from '@angular/router';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

  public user: User = {
    username: "",
    picture: "",
    nickname: "",
    updated_at: ""
  };

  constructor(private httpClient: HttpClient, private router: Router) { }

  ngOnInit() {
    this.httpClient.get(environment.gateway + "/auth/user").subscribe((data: User) => {
      this.user = data;
    });
  }

  onLogout() {
    window.location.href = environment.gateway + "/auth/logout";
  }
}

class User {
  username: string;
  picture: string;
  nickname: string;
  updated_at: string;
}
