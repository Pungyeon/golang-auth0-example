import { HomeComponent } from "./home/home.component";
import { RouterModule, Routes } from '@angular/router';
import { NgModule } from "@angular/core";
import { AuthGuardService } from "./service/auth-guard.service";
import { CallbackComponent } from "./callback/callback.component";
import { ChatComponent } from "src/app/chat/chat.component";

const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomeComponent },
    { path: 'chat', component: ChatComponent, canActivate: [AuthGuardService] },
    { path: 'callback', component: CallbackComponent }
  ];
  
  @NgModule({
    imports: [ RouterModule.forRoot(routes) ],
    exports: [ RouterModule ]
  })
  export class AppRoutingModule { }