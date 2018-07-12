import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule } from '@angular/common/http';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { AuthGuardService } from './service/auth-guard.service';
import { AuthService } from 'src/app/service/auth.service';
import { CallbackComponent } from 'src/app/callback/callback.component';
import { TodoComponent } from './todo/todo.component';
import { TodoService } from './service/todo.service';
import { FormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    CallbackComponent,
    TodoComponent
  ],
  imports: [
    NgbModule.forRoot(),
    AppRoutingModule,
    BrowserModule,
    FormsModule,
    HttpClientModule
  ],
  providers: [AuthGuardService, AuthService, TodoService],
  bootstrap: [AppComponent]
})
export class AppModule { }
