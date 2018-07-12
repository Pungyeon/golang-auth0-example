import { CanActivate, Router } from "@angular/router";
import { Injectable } from "@angular/core";
import { AuthService } from "./auth.service";
import { environment } from "src/environments/environment";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs";


@Injectable()
export class AuthGuardService implements CanActivate {
    constructor(private auth: AuthService, private router: Router) {}

    canActivate(): boolean {
        if (this.auth.isAuthenticated()) {
            return true;
        }

        this.auth.login()
        // this.router.navigate(['/auth/login']);
    }
}
