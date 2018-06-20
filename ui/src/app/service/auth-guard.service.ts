import { CanActivate, Router } from "@angular/router";
import { Injectable } from "@angular/core";
import { AuthService } from "./auth.service";
import { environment } from "src/environments/environment";


@Injectable()
export class AuthGuardService implements CanActivate {
    constructor(private auth: AuthService, private router: Router) {}

    canActivate(): boolean {
        if (this.auth.isAuthorized()) {
            this.router.navigate([environment + '/auth/login']);
            return false;
        }
        return true;
    }
}
