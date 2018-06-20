import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { HttpClient, HttpErrorResponse } from "@angular/common/http";

@Injectable()
export class AuthService {
    constructor(private httpClient: HttpClient)  {}

    isAuthorized() {
        this.httpClient.get(environment.gateway + "/auth/check-auth").subscribe((data) => {
            console.log(data);
            return true;
        }, (err) => {
            return false;
        })
    }

    AuthOK() {
        return this.httpClient.get(environment.gateway + "/auth/check-auth")
    }
}