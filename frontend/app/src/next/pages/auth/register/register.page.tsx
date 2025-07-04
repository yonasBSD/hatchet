import { Link, useNavigate } from 'react-router-dom';
import { UserRegisterForm } from './components/user-register-form';
import React from 'react';
import useApiMeta from '@/next/hooks/use-api-meta';
import {
  AuthLayout,
  GoogleLogin,
  GithubLogin,
  OrContinueWith,
} from '../components/shared-auth-components';
import useUser from '@/next/hooks/use-user';
import { ROUTES } from '@/next/lib/routes';
import { useTenantDetails } from '@/next/hooks/use-tenant';
import { Loading } from '@/components/ui/loading';

export default function Register() {
  const { oss, isLoading } = useApiMeta();

  if (isLoading) {
    return <Loading />;
  }

  const schemes = oss?.auth?.schemes || [];
  const basicEnabled = schemes.includes('basic');
  const googleEnabled = schemes.includes('google');
  const githubEnabled = schemes.includes('github');

  let prompt = 'Create an account to get started.';

  if (basicEnabled && (googleEnabled || githubEnabled)) {
    prompt =
      'Enter your email and password to create an account, or continue with a supported provider.';
  } else if (googleEnabled || githubEnabled) {
    prompt = 'Continue with a supported provider.';
  } else if (basicEnabled) {
    prompt = 'Create an account to get started.';
  } else {
    prompt = 'No login methods are enabled.';
  }

  const forms = [
    basicEnabled && <BasicRegister />,
    googleEnabled && <GoogleLogin />,
    githubEnabled && <GithubLogin />,
  ].filter(Boolean);

  return (
    <AuthLayout title="Create an account" prompt={prompt}>
      {forms.map((form, index) => (
        <React.Fragment key={index}>
          {form}
          {basicEnabled && schemes.length >= 2 && index === 0 ? (
            <OrContinueWith />
          ) : null}
        </React.Fragment>
      ))}
      <div className="flex flex-col space-y-2">
        <p className="text-sm text-gray-700 dark:text-gray-300">
          Already have an account?{' '}
          <Link
            to={ROUTES.auth.login}
            className="underline underline-offset-4 hover:text-primary"
          >
            Log in
          </Link>
        </p>
      </div>
    </AuthLayout>
  );
}

function BasicRegister() {
  const navigate = useNavigate();
  const { register, memberships } = useUser();
  const { defaultTenant } = useTenantDetails();

  const tenant = defaultTenant || memberships?.[0]?.tenant;

  return (
    <UserRegisterForm
      isLoading={register.isPending}
      onSubmit={async (data) => {
        const user = await register.mutateAsync(data);
        if (user && tenant?.metadata.id) {
          navigate(ROUTES.runs.list(tenant.metadata.id));
          return;
        }

        if (!tenant?.metadata.id) {
          navigate(ROUTES.onboarding.newTenant);
        }
      }}
      apiError={register.error?.message}
    />
  );
}
