import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  constructor(private router: Router) { }

  ngOnInit() {
  }

  onLogin() {
    const redirectUrl = environment.gateway + '/auth/login';
    if (environment.gateway === "") {
        this.router.navigate(['/auth/login']);   
    } else {
        window.location.href = redirectUrl;
    }
  }

  onUserDetails() {
    this.router.navigate(['/user']);   
  }
}
