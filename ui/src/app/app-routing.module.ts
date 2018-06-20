import { HomeComponent } from "./home/home.component";
import { RouterModule, Routes } from '@angular/router';
import { NgModule } from "@angular/core";
import { UserComponent } from "./user/user.component";
import { AuthGuardService } from "./service/auth-guard.service";

const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomeComponent },
    { path: 'user', component: UserComponent, canActivate: [AuthGuardService] }
  ];
  
  @NgModule({
    imports: [ RouterModule.forRoot(routes) ],
    exports: [ RouterModule ]
  })
  export class AppRoutingModule { }