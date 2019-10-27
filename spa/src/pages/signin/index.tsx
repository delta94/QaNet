import React from 'react';
import { NextPage } from 'next';
import Router from 'next/router';
import SignIn from '../../components/signin';
import { ICustomContext } from '../_app';

interface IProps {}

const SignInPage: NextPage<IProps> = function() {
  return <SignIn />;
};

interface Context extends ICustomContext {
  isLoggedIn: boolean;
}

SignInPage.getInitialProps = async function({ res, isLoggedIn }: Context) {
  if (isLoggedIn) {
    if (typeof window === 'undefined') {
      (res as any).redirect('/');
    } else {
      Router.push('/');
    }
  }

  return {};
};

export default SignInPage;