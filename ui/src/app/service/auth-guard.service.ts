import { CanActivate, Router } from "@angular/router";
import { Injectable } from "@angular/core";
import { AuthService } from "./auth.service";
import { environment } from "src/environments/environment";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs";


@Injectable()
export class AuthGuardService implements CanActivate {
    constructor(private auth: AuthService, private router: Router) {}

    canActivate(): Observable<boolean> {
        return new Observable<boolean>((observer) =>
        {
            this.auth.AuthOK().toPromise().then((data) => {
                observer.next(true);
                observer.complete();
            }).catch((err) => {
                const redirectUrl = environment.gateway + '/auth/login';
                if (environment.gateway === "") {
                    this.router.navigate(['/auth/login']);   
                } else {
                    window.location.href = redirectUrl;
                }
                observer.next(false);
                observer.complete();
            });
        });
    }
}
