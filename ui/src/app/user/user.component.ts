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

  public user;

  constructor(private httpClient: HttpClient, private router: Router) { }

  ngOnInit() {
    this.httpClient.get(environment.gateway + "/auth/user").subscribe((data) => {
      this.user = data;
    });
  }

  onLogout() {
    this.httpClient.get(environment.gateway + '/auth/logout').subscribe((data) => {
      console.log(data);
      this.router.navigate(['/home']);
    });
  }

}
