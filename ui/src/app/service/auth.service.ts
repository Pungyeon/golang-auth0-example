import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { HttpClient, HttpErrorResponse, HttpHeaders } from "@angular/common/http";
import { Router } from "@angular/router";

import * as auth0 from 'auth0-js';

(window as any).global = window;

@Injectable()
export class AuthService {
    constructor(public router: Router)  {}

    access_token: string;
    id_token: string;
    expires_at: string;

    auth0 = new auth0.WebAuth({
        clientID: 'bpF1FvreQgp1PIaSQm3fpCaI0A3TCz5T',
        domain: 'pungy.eu.auth0.com',
        responseType: 'token id_token',
        audience: 'https://pungy.eu.auth0.com/userinfo',
        redirectUri: environment.callback,
        scope: 'openid'
    });

    public login(): void {
        this.auth0.authorize();
    }

    // ...
    public handleAuthentication(): void {
        this.auth0.parseHash((err, authResult) => {
            if (authResult && authResult.accessToken && authResult.idToken) {
                window.location.hash = '';
                this.setSession(authResult);
                this.router.navigate(['/home']);
            } else if (err) {
                this.router.navigate(['/home']);
                console.log(err);
            }
        });
    }

    private setSession(authResult): void {
        // Set the time that the Access Token will expire at
        const expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
        this.access_token = authResult.accessToken;
        this.id_token = authResult.id_token;
        this.expires_at = authResult.expiresAt;
    }

    public logout(): void {
        this.access_token = null;
        this.id_token = null;
        this.expires_at = null;
        // Go back to the home route
        this.router.navigate(['/']);
    }

    public isAuthenticated(): boolean {
        // Check whether the current time is past the
        // Access Token's expiry time
        const expiresAt = JSON.parse(this.expires_at || '{}');
        return new Date().getTime() < expiresAt;
    }

    public createAuthHeaderValue(): string {
        if (this.id_token == "") {
            "";
        }
        return 'Bearer ' + this.id_token;
    }
}