import { Component, OnInit, ViewChild } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AuthService } from './auth.service';
import { TokenResponse } from './token-response.model';
import { StorageService } from '../storage.service';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-login',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent implements OnInit {
  @ViewChild('signinForm') signinForm: NgForm;

  errorMessage: string;
  isFetching = false;

  constructor(
    private activatedRoute: ActivatedRoute,
    private authService: AuthService,
    private storageService: StorageService,
    private router: Router
  ) {}

  ngOnInit() {}

  public login() {
    this.errorMessage = '';
    this.isFetching = true;

    const { userid, password } = this.signinForm.value;

    const successCb = (res: TokenResponse) => {
      const returnUrl = this.activatedRoute.snapshot.queryParams['returnUrl'] || '/questions';

      this.isFetching = false;
      this.storageService.saveToken(res.accessToken);

      this.storageService.save('userId', userid);

      this.authService.notifyLoggedIn();
      this.router.navigateByUrl(returnUrl);
    };

    const errorCb = error => {
      this.errorMessage = error;
      this.isFetching = false;
    };

    this.authService
      .login({ userid, password })
      .subscribe(successCb, errorCb);
  }

  onChange() {
    this.errorMessage = '';
  }
}
