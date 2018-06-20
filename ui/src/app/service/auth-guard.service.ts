import { CanActivate, Router } from "@angular/router";
import { Injectable } from "@angular/core";
import { AuthService } from "./auth.service";
import { environment } from "src/environments/environment";
import { HttpErrorResponse } from "@angular/common/http";


@Injectable()
export class AuthGuardService implements CanActivate {
    constructor(private auth: AuthService, private router: Router) {}

    canActivate(): boolean {
        if (!this.auth.isAuthorized()) {
            const redirectUrl = environment.gateway + '/auth/login';
            if (environment.gateway === "") {
                this.router.navigate(['/auth/login']);   
            } else {
                window.location.href = redirectUrl;
            }
            return false;
        } 
        return true;
    }
}
