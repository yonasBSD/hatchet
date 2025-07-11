import { Button } from '@/components/v1/ui/button';
import { cn } from '@/lib/utils';
import {
  BuildingOffice2Icon,
  // ChartBarSquareIcon,
  CheckIcon,
} from '@heroicons/react/24/outline';
import invariant from 'tiny-invariant';
import {
  Command,
  CommandEmpty,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '@/components/v1/ui/command';
import { Link } from 'react-router-dom';
import { TenantMember, TenantVersion } from '@/lib/api';
import { CaretSortIcon, PlusCircledIcon } from '@radix-ui/react-icons';
import {
  PopoverTrigger,
  Popover,
  PopoverContent,
} from '@radix-ui/react-popover';
import React from 'react';
import { Spinner } from '@/components/v1/ui/loading.tsx';
import useApiMeta from '@/pages/auth/hooks/use-api-meta';
import { useTenantDetails } from '@/hooks/use-tenant';

interface TenantSwitcherProps {
  className?: string;
  memberships: TenantMember[];
}
export function TenantSwitcher({
  className,
  memberships,
}: TenantSwitcherProps) {
  const meta = useApiMeta();
  const { setTenant: setCurrTenant, tenant: currTenant } = useTenantDetails();
  const [open, setOpen] = React.useState(false);

  if (!currTenant) {
    return <Spinner />;
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          aria-label="Select a team"
          className={cn('w-full justify-between', className)}
        >
          <BuildingOffice2Icon className="mr-2 h-4 w-4" />
          {currTenant.name}
          <CaretSortIcon className="ml-auto h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent side="right" className="w-full p-0 mb-6 z-50">
        <Command className="min-w-[260px]" value={currTenant.slug}>
          <CommandList>
            <CommandEmpty>No tenants found.</CommandEmpty>
            {memberships.map((membership) => (
              <CommandItem
                key={membership.metadata.id}
                onSelect={() => {
                  invariant(membership.tenant);
                  setCurrTenant(membership.tenant);
                  setOpen(false);

                  if (membership.tenant.version === TenantVersion.V0) {
                    // Hack to wait for next event loop tick so local storage is updated
                    setTimeout(() => {
                      window.location.href = `/workflow-runs?tenant=${membership.tenant?.metadata.id}`;
                    }, 0);
                  }
                }}
                value={membership.tenant?.slug}
                className="text-sm cursor-pointer"
              >
                <BuildingOffice2Icon className="mr-2 h-4 w-4" />
                {membership.tenant?.name}
                <CheckIcon
                  className={cn(
                    'ml-auto h-4 w-4',
                    currTenant.slug === membership.tenant?.slug
                      ? 'opacity-100'
                      : 'opacity-0',
                  )}
                />
              </CommandItem>
            ))}
          </CommandList>
          {meta.data?.allowCreateTenant && (
            <>
              <CommandSeparator />
              <CommandList>
                <Link to="/onboarding/create-tenant">
                  <CommandItem className="text-sm cursor-pointer">
                    <PlusCircledIcon className="mr-2 h-4 w-4" />
                    New Tenant
                  </CommandItem>
                </Link>
              </CommandList>
            </>
          )}
        </Command>
      </PopoverContent>
    </Popover>
  );
}
