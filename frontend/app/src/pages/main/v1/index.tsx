import { Button } from '@/components/v1/ui/button';
import { cn } from '@/lib/utils';
import {
  CalendarDaysIcon,
  CpuChipIcon,
  PlayIcon,
  ScaleIcon,
  ServerStackIcon,
  Squares2X2Icon,
} from '@heroicons/react/24/outline';

import {
  Link,
  Outlet,
  useLocation,
  useNavigate,
  useOutletContext,
} from 'react-router-dom';
import { Tenant, TenantMember, TenantUIVersion } from '@/lib/api';
import { ClockIcon, GearIcon } from '@radix-ui/react-icons';
import React, { useCallback, useEffect } from 'react';
import {
  MembershipsContextType,
  UserContextType,
  useContextFromParent,
} from '@/lib/outlet';
import { useTenant } from '@/lib/atoms';
import { Loading } from '@/components/v1/ui/loading.tsx';
import { TenantSwitcher } from '@/components/v1/molecules/nav-bar/tenant-switcher';
import {
  Collapsible,
  CollapsibleContent,
} from '@/components/v1/ui/collapsible';
import useCloudApiMeta from '@/pages/auth/hooks/use-cloud-api-meta';
import useCloudFeatureFlags from '@/pages/auth/hooks/use-cloud-feature-flags';
import { useSidebar } from '@/components/sidebar-provider';
import invariant from 'tiny-invariant';
import { ROUTES } from '@/next/lib/routes';
import { SquareActivityIcon } from 'lucide-react';

function Main() {
  const ctx = useOutletContext<UserContextType & MembershipsContextType>();
  const navigate = useNavigate();
  const { user, memberships } = ctx;

  const { tenant: currTenant } = useTenant();

  const childCtx = useContextFromParent({
    user,
    memberships,
    tenant: currTenant,
  });

  useEffect(() => {
    if (!currTenant) {
      return;
    }

    if (currTenant.uiVersion === TenantUIVersion.V1) {
      // Hard redirect here because the navigate hook is racy with other url updates
      window.location.href = ROUTES.runs.list(currTenant.metadata.id);
    }
  }, [currTenant, navigate]);

  if (!user || !memberships || !currTenant) {
    return <Loading />;
  }

  return (
    <div className="flex flex-row flex-1 w-full h-full">
      <Sidebar memberships={memberships} currTenant={currTenant} />
      <div className="p-8 flex-grow overflow-y-auto overflow-x-hidden">
        <Outlet context={childCtx} />
      </div>
    </div>
  );
}

export default Main;

interface SidebarProps extends React.HTMLAttributes<HTMLDivElement> {
  memberships: TenantMember[];
  currTenant: Tenant;
}

function Sidebar({ className, memberships, currTenant }: SidebarProps) {
  const { sidebarOpen, setSidebarOpen } = useSidebar();
  const { tenant } = useTenant();

  invariant(tenant);

  const meta = useCloudApiMeta();
  const featureFlags = useCloudFeatureFlags(currTenant.metadata.id);

  const onNavLinkClick = useCallback(() => {
    if (window.innerWidth > 768) {
      return;
    }

    setSidebarOpen('closed');
  }, [setSidebarOpen]);

  if (sidebarOpen === 'closed') {
    return null;
  }

  const workers = [
    <SidebarButtonSecondary
      key={1}
      onNavLinkClick={onNavLinkClick}
      to="/v1/workers/all"
      name="All Workers"
    />,
    <SidebarButtonSecondary
      key={1}
      onNavLinkClick={onNavLinkClick}
      to="/v1/workers/webhook"
      name="Webhook Workers"
    />,
  ];

  return (
    <div
      className={cn(
        'h-full border-r overflow-y-auto w-full md:min-w-80 md:w-80 top-16 absolute z-[100] md:relative md:top-0 md:bg-[unset] md:dark:bg-[unset] bg-slate-100 dark:bg-slate-900',
        className,
      )}
    >
      <div className="flex flex-col justify-between items-start space-y-4 px-4 py-4 h-full pb-16 md:pb-4">
        <div className="grow w-full">
          <div className="py-2">
            <h2 className="mb-2 text-lg font-semibold tracking-tight">
              Activity
            </h2>
            <div className="flex flex-col gap-y-1">
              <SidebarButtonPrimary
                key={1}
                onNavLinkClick={onNavLinkClick}
                to="/v1/runs"
                name="Runs"
                icon={<PlayIcon className="mr-2 h-4 w-4" />}
              />
              <SidebarButtonPrimary
                key={2}
                onNavLinkClick={onNavLinkClick}
                to="/v1/events"
                name="Events"
                icon={<SquareActivityIcon className="mr-2 h-4 w-4" />}
              />
            </div>
          </div>
          <div className="py-2">
            <h2 className="mb-2 text-lg font-semibold tracking-tight">
              Triggers
            </h2>
            <div className="space-y-1">
              <SidebarButtonPrimary
                key={3}
                onNavLinkClick={onNavLinkClick}
                to="/v1/scheduled"
                name="Scheduled Runs"
                icon={<CalendarDaysIcon className="mr-2 h-4 w-4" />}
              />
              <SidebarButtonPrimary
                key={4}
                onNavLinkClick={onNavLinkClick}
                to="/v1/cron-jobs"
                name="Cron Jobs"
                icon={<ClockIcon className="mr-2 h-4 w-4" />}
              />
              {/* <SidebarButtonPrimary
                key={5}
                onNavLinkClick={onNavLinkClick}
                to="/v1/events"
                name="Events"
                icon={<QueueListIcon className="mr-2 h-4 w-4" />}
              /> */}
            </div>
          </div>
          <div className="py-2">
            <h2 className="mb-2 text-lg font-semibold tracking-tight">
              Resources
            </h2>
            <div className="space-y-1">
              <SidebarButtonPrimary
                key={1}
                onNavLinkClick={onNavLinkClick}
                to="/v1/tasks"
                name="Tasks & Workflows"
                icon={<Squares2X2Icon className="mr-2 h-4 w-4" />}
              />
              <SidebarButtonPrimary
                key={2}
                onNavLinkClick={onNavLinkClick}
                to="/v1/workers/all"
                name="Workers"
                icon={<ServerStackIcon className="mr-2 h-4 w-4" />}
                prefix="/v1/workers"
                collapsibleChildren={workers}
              />
              {featureFlags?.data['managed-worker'] && (
                <SidebarButtonPrimary
                  key={3}
                  onNavLinkClick={onNavLinkClick}
                  to="/v1/managed-workers"
                  name="Managed Compute"
                  icon={<CpuChipIcon className="mr-2 h-4 w-4" />}
                />
              )}
              <SidebarButtonPrimary
                key={4}
                onNavLinkClick={onNavLinkClick}
                to="/v1/rate-limits"
                name="Rate Limits"
                icon={<ScaleIcon className="mr-2 h-4 w-4" />}
              />
            </div>
          </div>
          <div className="py-2">
            <h2 className="mb-2 text-lg font-semibold tracking-tight">
              Settings
            </h2>
            <div className="space-y-1">
              <SidebarButtonPrimary
                key={1}
                onNavLinkClick={onNavLinkClick}
                to="/v1/tenant-settings/overview"
                prefix="/v1/tenant-settings"
                name="General"
                icon={<GearIcon className="mr-2 h-4 w-4" />}
                collapsibleChildren={[
                  <SidebarButtonSecondary
                    key={1}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/overview"
                    name="Overview"
                  />,
                  <SidebarButtonSecondary
                    key={2}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/api-tokens"
                    name="API Tokens"
                  />,
                  <SidebarButtonSecondary
                    key={3}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/github"
                    name="Github"
                  />,
                  <SidebarButtonSecondary
                    key={5}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/members"
                    name="Members"
                  />,
                  <SidebarButtonSecondary
                    key={6}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/billing-and-limits"
                    name={
                      meta?.data.canBill
                        ? 'Billing & Limits'
                        : 'Resource Limits'
                    }
                  />,
                  <SidebarButtonSecondary
                    key={7}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/alerting"
                    name="Alerting"
                  />,
                  <SidebarButtonSecondary
                    key={8}
                    onNavLinkClick={onNavLinkClick}
                    to="/v1/tenant-settings/ingestors"
                    name="Ingestors"
                  />,
                ]}
              />
            </div>
          </div>
        </div>
        <TenantSwitcher memberships={memberships} currTenant={currTenant} />
      </div>
    </div>
  );
}

function SidebarButtonPrimary({
  onNavLinkClick,
  to,
  name,
  icon,
  prefix,
  collapsibleChildren = [],
}: {
  onNavLinkClick: () => void;
  to: string;
  name: string;
  icon: React.ReactNode;
  prefix?: string;
  collapsibleChildren?: React.ReactNode[];
}) {
  const location = useLocation();
  const open = location.pathname.startsWith(prefix || to);
  const selected = !prefix && location.pathname === to;

  const primaryLink = (
    <Link to={to} onClick={onNavLinkClick}>
      <Button
        variant="ghost"
        className={cn(
          'w-full justify-start pl-2',
          selected && 'bg-slate-200 dark:bg-slate-800',
        )}
      >
        {icon}
        {name}
      </Button>
    </Link>
  );

  return collapsibleChildren.length == 0 ? (
    primaryLink
  ) : (
    <Collapsible
      open={open}
      // onOpenChange={setIsOpen}
      className="w-full"
    >
      {primaryLink}
      <CollapsibleContent className={'space-y-2 ml-4 border-l border-muted'}>
        {collapsibleChildren}
      </CollapsibleContent>
    </Collapsible>
  );
}

function SidebarButtonSecondary({
  onNavLinkClick,
  to,
  name,
  prefix,
}: {
  onNavLinkClick: () => void;
  to: string;
  name: string;
  prefix?: string;
}) {
  const location = useLocation();
  const hasPrefix = prefix && location.pathname.startsWith(prefix);
  const selected = hasPrefix || location.pathname === to;

  return (
    <Link to={to} onClick={onNavLinkClick}>
      <Button
        variant="ghost"
        size="sm"
        className={cn(
          'w-[calc(100%-3px)] justify-start pl-3 pr-0 ml-1 mr-3 my-[1px]',
          selected && 'bg-slate-200 dark:bg-slate-800',
        )}
      >
        {name}
      </Button>
    </Link>
  );
}
